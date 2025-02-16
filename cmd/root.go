/*
Copyright © 2025 Francesco Orlando scrotadamus@insiberia.net
*/
package cmd

import (
	"os"

	"github.com/scrotadamus/ghligh/cmd/tag"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghligh",
	Short: "pdf highlights swiss knife",
	Long:  `ghligh can be used to manipulate pdf files in various ways.`,

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(tag.TagCmd)
}
