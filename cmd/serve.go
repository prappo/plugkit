/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

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

		pluginDir := args[0]
		// Run system command 'npm run dev:server' in the plugin directory using cobra

		sysCmd := exec.Command("npm", "--prefix", pluginDir, "run", "dev:server")
		sysCmd.Stdout = os.Stdout
		sysCmd.Stderr = os.Stderr
		err := sysCmd.Run()
		if err != nil {
			fmt.Println("Error running npm run dev:server", err)
			return
		}
		fmt.Println("npm run dev:server output:")

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
