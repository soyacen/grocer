package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/yaml"

	"github.com/cockroachdb/errors"
	"github.com/soyacen/grocer/internal/edit"
	"github.com/spf13/cobra"
)

// cronjobCmd represents the cronjob command
var cronjobCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "add cronjob",
	RunE:  cronjobRun,
}

type cronjobFlags struct {
	Name string
	Dir  string
}

var cronjobFlag cronjobFlags

func init() {
	rootCmd.AddCommand(cronjobCmd)
	cronjobCmd.Flags().StringVarP(&cronjobFlag.Name, "name", "n", "", "cron job name")
	_ = cronjobCmd.MarkFlagRequired("name")
	cronjobCmd.Flags().StringVarP(&cronjobFlag.Dir, "dir", "d", "", "project directory, default is current directory")
}

func cronjobRun(_ *cobra.Command, _ []string) error {
	srcMod, srcModVers, err := getSrcModInfo()
	if err != nil {
		return err
	}

	info, err := getGoModInfo(srcMod, srcModVers)
	if err != nil {
		return err
	}

	dir, err := getProjectDir(cronjobFlag.Dir, "")
	if err != nil {
		return err
	}

	// Dir must exist and must be non-empty.
	de, err := os.ReadDir(dir)
	if err != nil || len(de) == 0 {
		return errors.New("target directory does not exist or is empty")
	}

	dstMod, err := readMod(dir)
	if err != nil {
		return err
	}

	// Copy from module cache into new directory, making edits as needed.
	if err := filepath.WalkDir(info.Dir, func(src string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		rel, err := filepath.Rel(info.Dir, src)
		if err != nil {
			return errors.WithStack(err)
		}

		prefixs := []string{
			"cmd/cronjob",
			"deploy/cronjob",
			"internal/cronjob",
		}
		for _, prefix := range prefixs {
			if !strings.HasPrefix(rel, prefix) {
				return nil
			}
		}

		dst := filepath.Join(dir, rel)
		if d.IsDir() {
			if err := os.MkdirAll(dst, 0o777); err != nil {
				return errors.WithStack(err)
			}
			return nil
		}

		data, err := os.ReadFile(src)
		if err != nil {
			return errors.WithStack(err)
		}

		switch rel {
		case "cmd/cronjob.go":
			data, err = fixCmdCronjobGo(data, dir)
		case "deploy/cronjob.yaml":
			data, err = fixDeployCronjobYaml(data, dir, dstMod)
		case "internal/cronjob/wire.go":
		case "internal/cronjob/service.go":
		case "internal/cronjob/repo.go":
		case "internal/cronjob/repo.go":

			data, err = fixRepo(data, dir, dstMod)
		}
		if err != nil {
			return err
		}
		if strings.HasSuffix(rel, ".go") {
			isRoot := !strings.Contains(rel, string(filepath.Separator))
			data, err = fixGo(data, rel, srcMod, dstMod, isRoot)
			if err != nil {
				return err
			}
		}
		if err := os.WriteFile(dst, data, 0o666); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		return err
	}

	log.Printf("add cron job %s in %s", cronjobFlag.Name, dir)
	return nil
}

func fixDeployCronjobYaml(data []byte, dir string, mod string) ([]byte, error) {
	// 解析YAML内容到CronJob结构体
	var cronJob batchv1.CronJob
	if err := yaml.Unmarshal(data, &cronJob); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal cronjob yaml")
	}

	// 从mod参数中提取项目名称作为应用名称
	appName := path.Base(mod)

	// 更新metadata中的name
	cronJob.Name = strings.Replace(cronJob.Name, "grocer", appName, -1)

	// 如果Labels存在，更新app标签
	if cronJob.Labels != nil {
		cronJob.Labels["app"] = appName
	} else {
		cronJob.Labels = map[string]string{
			"app": appName,
		}
	}

	// 如果Annotations存在，更新描述
	if cronJob.Annotations != nil {
		// 使用新的描述格式，包含schedule信息
		schedule := cronJob.Spec.Schedule
		cronJob.Annotations["description"] = fmt.Sprintf("定时任务%s，%s", cronJob.Name, schedule)
	} else {
		// 使用新的描述格式
		schedule := cronJob.Spec.Schedule
		cronJob.Annotations = map[string]string{
			"description": fmt.Sprintf("定时任务%s，%s", cronJob.Name, schedule),
		}
	}

	// 更新spec部分中的container名称和镜像
	jobTemplate := &cronJob.Spec.JobTemplate
	container := &jobTemplate.Spec.Template.Spec.Containers[0]

	// 更新容器名称
	container.Name = strings.Replace(container.Name, "grocer", appName, -1)

	// 从mod中提取镜像名称，基于项目名称
	imageName := strings.Replace(mod, "/", "-", -1)    // 将路径分隔符替换为短横线
	image := fmt.Sprintf("%s:%s", imageName, "latest") // 使用latest作为默认版本
	container.Image = image

	// 生成修改后的YAML
	result, err := yaml.Marshal(&cronJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal cronjob yaml")
	}

	return result, nil
}

func fixCmdCronjobGo(data []byte, dir string) ([]byte, error) {
	filename := "cmd/cronjob.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, data, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", filename)
	}

	buf := edit.NewBuffer(data)

	// 遍历 AST 查找 rootCmd 变量的 Use 字段
	ast.Inspect(f, func(n ast.Node) bool {
		x, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		sel, ok := x.Type.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "cobra" || sel.Sel.Name != "Command" {
			return true
		}

		// 遍历结构体字段
		for _, elt := range x.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != "Use" {
				continue
			}

			val, ok := kv.Value.(*ast.BasicLit)
			if !ok || val.Kind != token.STRING {
				continue
			}

			oldVal, _ := strconv.Unquote(val.Value)
			newVal := strings.Replace(oldVal, "cronjob", path.Base(dir), -1)
			if newVal != oldVal {
				buf.Replace(
					fset.Position(kv.Value.Pos()).Offset,
					fset.Position(kv.Value.End()).Offset,
					strconv.Quote(newVal),
				)
			}
		}
		return true
	})

	return buf.Bytes(), nil
}
