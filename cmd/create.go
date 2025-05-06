/*
Copyright © 2025 Prappo <prappo.prince@gmail.com>
*/
package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

		if len(args) == 0 {
			cmd.Help()
			return
		}

		// If more than one argument is provided, show the error message
		if len(args) > 1 {
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
	// Download the plugin boilerplate from the GitHub repository ( https://github.com/prappo/wordpress-plugin-boilerplate )

	downloadURL := "https://github.com/prappo/wordpress-plugin-boilerplate/archive/refs/heads/main.zip"

	// Send GET request to the download URL

	response, err := http.Get(downloadURL)

	if err != nil {
		fmt.Println("Error downloading the plugin boilerplate")
	}

	defer response.Body.Close()

	out, err := os.Create(pluginName + ".zip")

	if err != nil {
		fmt.Println("Error creating the plugin boilerplate file")
	}

	// Copy response body to the file
	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println("Error copying the plugin boilerplate to the file")
	}

	// Close the output file
	out.Close()

	// Unzip the plugin boilerplate
	fmt.Println("Unzipping the plugin boilerplate...")

	zipReader, err := zip.OpenReader(pluginName + ".zip")
	if err != nil {
		fmt.Println("Error opening the zip file")
	}

	defer zipReader.Close()

	destinationDir := "./" + pluginName
	for _, file := range zipReader.File {
		// Skip the root directory entry
		if file.Name == "wordpress-plugin-boilerplate-main/" {
			continue
		}

		// Remove the boilerplate directory prefix from the path
		relativePath := file.Name
		if len(file.Name) > len("wordpress-plugin-boilerplate-main/") {
			relativePath = file.Name[len("wordpress-plugin-boilerplate-main/"):]
		}

		fmt.Println("Extracting file: ", relativePath)

		fpath := filepath.Join(destinationDir, relativePath)

		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Create the file
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				fmt.Println(err)
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				fmt.Println(err)
			}

			rc, err := file.Open()
			if err != nil {
				fmt.Println(err)
			}

			_, err = io.Copy(outFile, rc)
			if err != nil {
				fmt.Println(err)
			}

			outFile.Close()
			rc.Close()
		}
	}

	// Close the zip reader before removing the file
	zipReader.Close()

	// Remove the zip file
	removeZipFile(pluginName)

	fmt.Println("Plugin boilerplate created successfully")
}

func removeZipFile(pluginName string) {
	// Add a small delay to ensure all file handles are released
	time.Sleep(100 * time.Millisecond)

	err := os.Remove("./" + pluginName + ".zip")
	if err != nil {
		fmt.Println("Error removing the zip file")
		fmt.Println(err)
		return
	}

	fmt.Println("Zip file removed successfully")
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
