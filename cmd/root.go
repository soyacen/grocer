package cmd

import (
	"os"
	"regexp"

	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
)

type flags struct {
	Name string
	Dir  string
}

func (f flags) IsValid() error {
	// 编译正则表达式：以字母开头，后跟字母、数字或下划线
	validName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !validName.MatchString(f.Name) {
		return errors.New("name must consist of alphanumeric characters and underscores, and start with a letter")
	}
	return nil
}

var flag flags

var rootCmd = &cobra.Command{
	Use:          "grocer",
	Short:        "",
	Long:         ``,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
