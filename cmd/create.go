/*
Copyright © 2025 Prappo <prappo.prince@gmail.com>

*/
package cmd

import (
	"fmt"

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
		// If not plugin name is provided, show the help message

		if(len(args) == 0) {
			cmd.Help()
			return
		}

		// If more than one argument is provided, show the error message
		if ( len(args) > 1 ) {
			fmt.Println("Invalid plugin name provided")
			fmt.Println("Correct name should not contain spaces or special characters")
			fmt.Println("Example:")
			fmt.Println("plugkit create my-plugin")
			return
		}

		pluginName := args[0]

		downloadPluginBoilerplate(pluginName)
	},
}

func downloadPluginBoilerplate(pluginName string) {
	// GitHub URL to the plugin boilerplate
	fmt.Println("Downloading plugin boilerplate...")
	fmt.Println(pluginName)


}	

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
