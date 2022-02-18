/*
Copyright © 2022 Joshua Leong <juushdev@gmail.com>

*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes spark executable from svm directory",
	Long: `Removes spark executable from svm directory. For example:

svm remove 2.2.2-with-hadoop-2.7`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateRemoveArgs(args); err != nil {
			return err
		}
		fmt.Println("remove called")
		return nil
	},
}

func validateRemoveArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("No version specified. Use `svm list` to view all installed versions \n")
	}
	if len(args) > 1 {
		return errors.New("Multiple versions declared. Please only specify one version \n")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
