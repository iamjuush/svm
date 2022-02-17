/*
Copyright © 2022 Joshua Leong <juushdev@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows all existing spark versions",
	Long: `Shows all existing spark versions. For example:

Example:
> svm list
3.1.2
2.2.2-with-hadoop-2.7.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		svmFilepath := filepath.Join(dirname, ".svm")
		fmt.Printf("Available spark versions:\n")
		svmDir, err := os.ReadDir(svmFilepath)
		if err != nil {
			return err
		}
		for _, dir := range svmDir {
			if dir.Name() == "bin" {
				continue
			}
			fmt.Println(dir.Name())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
