package cmd

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"golang.org/x/mod/module"
)

func getSrcModInfo() (string, string, error) {
	srcMod := "github.com/soyacen/grocer/internal/layout"
	srcModVers := srcMod + "@latest"
	srcMod, _, _ = strings.Cut(srcMod, "@")
	if err := module.CheckPath(srcMod); err != nil {
		return "", "", errors.Wrap(err, "invalid source module name")
	}
	return srcMod, srcModVers, nil
}

type ModInfo struct {
	Dir string
}

func getGoModInfo(srcMod, srcModVers string) (*ModInfo, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "mod", "download", "-json", srcModVers)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.Errorf("go mod download -json %s: %v\n%s%s", srcModVers, err, stderr.Bytes(), stdout.Bytes())
	}
	var info ModInfo
	if err := json.Unmarshal(stdout.Bytes(), &info); err != nil {
		return nil, errors.Errorf("go mod download -json %s: invalid JSON output: %v\n%s%s", srcMod, err, stderr.Bytes(), stdout.Bytes())
	}
	return &info, nil
}

func getProjectDir(dir, dstMod string) (string, error) {
	if dir == "" {
		dir = "."
		if dstMod != "" {
			dir += string(filepath.Separator) + path.Base(dstMod)
		}
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", errors.Wrap(err, "failed to get absolute path for target directory")
	}
	return absDir, nil
}

func readMod(dir string) (string, error) {
}
