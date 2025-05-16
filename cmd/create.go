/*
Copyright © 2025 Prappo <prappo.prince@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/prappo/plugkit/internal/commands"
	"github.com/prappo/plugkit/internal/config"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new WordPress plugin",
	Long: `Create a new WordPress plugin with the given name.

	Example:
	plugkit create my-plugin
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Error: Please provide a plugin folder name")
			fmt.Println("Example: plugkit create my-plugin")
			cmd.Help()
			return
		}

		folderName := args[0]
		pluginConfig, err := config.CollectPluginConfig(folderName)
		if err != nil {
			fmt.Printf("Error collecting plugin configuration: %v\n", err)
			os.Exit(1)
		}

		if err := commands.CreatePlugin(pluginConfig); err != nil {
			fmt.Printf("Error creating plugin: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
