package cmd

import (
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
		}
		if err != nil {
			return err
		}
		isRoot := !strings.Contains(rel, string(filepath.Separator))
		if strings.HasSuffix(rel, ".go") {
			data = fixGo(data, rel, srcMod, dstMod, isRoot)
		}
		if err := os.WriteFile(dst, data, 0o666); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}); err != nil {
		return err
	}

	log.Printf("add cron job %s in %s", dstMod, dir)
	return nil
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
