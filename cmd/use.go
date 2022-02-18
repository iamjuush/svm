/*
Copyright © 2022 Joshua Leong <juushdev@gmail.com>

*/
package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Sets the desired spark version to use",
	Long: `Sets the desired spark version to use. For example:

svm use 2.2.2-with-hadoop-2.7
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateUseArgs(args); err != nil {
			return err
		}
		dirname, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		sourceDir := filepath.Join(dirname, ".svm", args[0])
		targetDir := filepath.Join(dirname, ".svm", "active")
		_ = os.Remove(targetDir)
		if err = os.Symlink(sourceDir, targetDir); err != nil {
			return err
		}
		return nil
	},
}

func validateUseArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("No version specified. Use `svm list` to view all installed versions \n")
	}
	if len(args) > 1 {
		return errors.New("Multiple versions declared. Please only specify one version \n")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
