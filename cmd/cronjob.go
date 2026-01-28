package cmd

import (
	"bytes"
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
			"deploy/values/cronjob",
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
		case "deploy/values/cronjob.yaml":
			data = bytes.ReplaceAll(data, []byte("grocer-cronjob"), []byte(path.Base(dstMod)+"-"+cronjobFlag.Name))
		case "internal/cronjob/fx.go":
			data, err = fixInternalCronjobFxGo(data, dir, dstMod)
		case "internal/cronjob/repo.go":
			data, err = fixInternalCronjobRepoGo(data, dir, dstMod)
		case "internal/cronjob/repository.go":
			data, err = fixInternalCronjobRepositoryGo(data, dir, dstMod)
		case "internal/cronjob/service.go":
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

func fixInternalCronjobFxGo(data []byte, dir, mod string) ([]byte, error) {
	filename := "internal/cronjob/fx.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, data, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", filename)
	}

	buf := edit.NewBuffer(data)
	newName := cronjobFlag.Name

	// 遍历 AST 查找并替换所有相关的 cronjob 标识符
	ast.Inspect(f, func(n ast.Node) bool {
		// 处理 fx.Module 的参数，如 fx.Module("cronjob", ...)
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "fx" || sel.Sel.Name != "Module" {
			return true
		}
		if len(call.Args) == 0 {
			return true
		}
		arg, ok := call.Args[0].(*ast.BasicLit)
		if !ok || arg.Kind != token.STRING {
			return true
		}
		oldVal, _ := strconv.Unquote(arg.Value)
		if oldVal != "cronjob" {
			return true
		}
		buf.Replace(
			fset.Position(arg.Pos()).Offset,
			fset.Position(arg.End()).Offset,
			strconv.Quote(newName),
		)
		return true
	})

	return buf.Bytes(), nil
}

func fixInternalCronjobRepositoryGo(data []byte, dir, mod string) ([]byte, error) {
	filename := "internal/cronjob/repository.go"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, data, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s", filename)
	}

	buf := edit.NewBuffer(data)
	newName := cronjobFlag.Name

	// 遍历 AST 查找并替换所有相关的 cronjob 标识符
	ast.Inspect(f, func(n ast.Node) bool {
		// 处理 db, err, _ := dbs.Load("cronjob") 和 rd, err, _ := rds.Load("cronjob")
		// 替换调用参数中的 "cronjob" 为新名称，但不替换包名
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		ident, ok := sel.X.(*ast.Ident)
		if !ok || (ident.Name != "dbs" && ident.Name != "rds") || sel.Sel.Name != "Load" {
			return true
		}
		if len(call.Args) == 0 {
			return true
		}
		arg, ok := call.Args[0].(*ast.BasicLit)
		if !ok || arg.Kind != token.STRING {
			return true
		}
		oldVal, _ := strconv.Unquote(arg.Value)
		if oldVal != "cronjob" {
			return true
		}
		buf.Replace(
			fset.Position(arg.Pos()).Offset,
			fset.Position(arg.End()).Offset,
			strconv.Quote(newName),
		)
		return true
	})

	return buf.Bytes(), nil
}

func fixInternalCronjobRepoGo(data []byte, dir, mod string) ([]byte, error) {
	// 目前 repo.go 与 repository.go 处理逻辑相同
	return fixInternalCronjobRepositoryGo(data, dir, mod)
}

func fixCmdCronjobGo(data []byte, name string) ([]byte, error) {
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
			newVal := strings.Replace(oldVal, "cronjob", name, -1)
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
