/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the plugin in development mode",
	Long: `Run the plugin in development mode.

	Example:
	plugkit serve my-plugin
	OR
	plugkit serve
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Invalid number of arguments")
			cmd.Help()
			return
		}
		fmt.Println("serve called")
		fmt.Println(len(args))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
