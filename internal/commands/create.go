package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/prappo/plugkit/internal/config"
	"github.com/schollz/progressbar/v3"
)

func CreatePlugin(pluginConfig *config.PluginConfig) error {
	fmt.Printf("Creating WordPress plugin: %s\n", pluginConfig.PluginName)

	// Download boilerplate
	zipPath := pluginConfig.OriginalName + ".zip"
	if err := downloadBoilerplate(zipPath); err != nil {
		return fmt.Errorf("failed to download boilerplate: %w", err)
	}
	defer os.Remove(zipPath)

	// Extract files
	if err := extractFiles(zipPath, pluginConfig.OriginalName); err != nil {
		return fmt.Errorf("failed to extract files: %w", err)
	}

	// Update plugin files with configuration
	if err := updatePluginFiles(pluginConfig); err != nil {
		return fmt.Errorf("failed to update plugin files: %w", err)
	}

	fmt.Printf("\nPlugin '%s' created successfully! 🎉\n", pluginConfig.PluginName)
	cleanUp(pluginConfig.OriginalName)
	return nil
}

func downloadBoilerplate(zipPath string) error {
	fmt.Println("Downloading plugin boilerplate...")

	resp, err := http.Get(config.GetConfig()["download_url"])
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create progress bar
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	out, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a proxy reader that updates the progress bar
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	return err
}

func extractFiles(zipPath, pluginName string) error {
	fmt.Println("Extracting files...")

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create progress bar for extraction
	bar := progressbar.NewOptions(
		len(reader.File),
		progressbar.OptionSetDescription("Extracting"),
		progressbar.OptionSetWidth(15),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	prefix := "wordpress-plugin-boilerplate-main/"
	for _, file := range reader.File {
		// Skip root directory
		if file.Name == prefix {
			continue
		}

		// Get relative path
		relPath := file.Name
		if len(file.Name) > len(prefix) {
			relPath = file.Name[len(prefix):]
		}

		destPath := filepath.Join(pluginName, relPath)
		if err := extractFile(file, destPath); err != nil {
			return fmt.Errorf("failed to extract %s: %w", relPath, err)
		}

		bar.Add(1)
	}
	return nil
}

func extractFile(file *zip.File, destPath string) error {
	if file.FileInfo().IsDir() {
		return os.MkdirAll(destPath, os.ModePerm)
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
		return err
	}

	// Open source file
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Create destination file
	out, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy contents
	_, err = io.Copy(out, rc)
	return err
}

func updatePluginFiles(config *config.PluginConfig) error {
	// First rename the main plugin file
	oldMainFile := filepath.Join(config.OriginalName, "wordpress-plugin-boilerplate.php")
	newMainFile := filepath.Join(config.OriginalName, config.PluginFileName)
	if err := os.Rename(oldMainFile, newMainFile); err != nil {
		return fmt.Errorf("failed to rename main plugin file: %w", err)
	}

	// Define the replacements to make
	replacements := []struct {
		pattern     string
		replacement string
		paths       []string
		recursive   bool
	}{
		// Main plugin file replacements
		{"WordPress Plugin Boilerplate", config.PluginName, []string{config.PluginFileName}, false},
		{"A boilerplate for WordPress plugins.", config.PluginDescription, []string{config.PluginFileName}, false},
		{"Prappo", config.AuthorName, []string{config.PluginFileName}, false},
		{"https://prappo.github.io", config.AuthorURI, []string{config.PluginFileName}, false},
		{`Version: [0-9]+\.[0-9]+\.[0-9]+`, fmt.Sprintf("Version: %s", config.PluginVersion), []string{config.PluginFileName}, false},
		{"Text Domain: wordpress-plugin-boilerplate", fmt.Sprintf("Text Domain: %s", config.TextDomain), []string{config.PluginFileName}, false},

		// Namespace replacements
		{"namespace WordPressPluginBoilerplate", fmt.Sprintf("namespace %s", config.Namespace), []string{config.PluginFileName, "plugin.php"}, false},
		{"use WordPressPluginBoilerplate", fmt.Sprintf("use %s", config.Namespace), []string{config.PluginFileName, "plugin.php"}, false},
		{"namespace WordPressPluginBoilerplate", fmt.Sprintf("namespace %s", config.Namespace), []string{"includes"}, true},
		{"use WordPressPluginBoilerplate", fmt.Sprintf("use %s", config.Namespace), []string{"includes"}, true},
		{"WordPressPluginBoilerplate", config.Namespace, []string{"includes"}, true},

		// Database namespace replacements
		{"namespace WordPressPluginBoilerplate", fmt.Sprintf("namespace %s", config.Namespace), []string{"database"}, true},
		{"WordPressPluginBoilerplate", config.Namespace, []string{"database"}, true},
		{"use WordPressPluginBoilerplate", fmt.Sprintf("use %s", config.Namespace), []string{"database"}, true},

		// Libs namespace replacements
		{"namespace WordPressPluginBoilerplate", fmt.Sprintf("namespace %s", config.Namespace), []string{"libs"}, true},
		{"use WordPressPluginBoilerplate", fmt.Sprintf("use %s", config.Namespace), []string{"libs"}, true},

		// Function and constant replacements
		{"wordpress_plugin_boilerplate_", fmt.Sprintf("%s_", config.PluginPrefix), []string{"includes"}, true},
		{"WordPressPluginBoilerplate", config.MainClassName, []string{config.PluginFileName, "plugin.php"}, false},
		{"wordpress_plugin_boilerplate_init", config.MainFunctionName, []string{config.PluginFileName}, false},
		{"WORDPRESS_PLUGIN_BOILERPLATE_", fmt.Sprintf("%s_", config.ConstantPrefix), []string{config.PluginFileName, "includes", "plugin.php"}, true},
	}

	// Apply all replacements
	for _, r := range replacements {
		for _, path := range r.paths {
			fullPath := filepath.Join(config.OriginalName, path)
			if err := applyReplacements(fullPath, r.pattern, r.replacement, r.recursive); err != nil {
				return fmt.Errorf("failed to apply replacements in %s: %w", path, err)
			}
		}
	}

	return nil
}

func applyReplacements(path, pattern, replacement string, recursive bool) error {
	// If path is a directory and recursive is true, process all files in the directory
	if recursive {
		return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(filePath, ".php") {
				return replaceInFile(filePath, pattern, replacement)
			}
			return nil
		})
	}

	// If path is a file, process it directly
	return replaceInFile(path, pattern, replacement)
}

func replaceInFile(filePath, pattern, replacement string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Create a regex pattern
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %w", err)
	}

	// Replace the pattern
	newContent := re.ReplaceAll(content, []byte(replacement))

	return os.WriteFile(filePath, newContent, 0644)
}

func cleanUp(pluginName string) {
	// List of files and directories to remove
	filesToRemove := []string{
		"npm",
		".storybook",
		"documentation",
		".github",
	}

	// Remove files and directories
	for _, file := range filesToRemove {
		fpath := filepath.Join(pluginName, file)
		if _, err := os.Stat(fpath); err == nil {
			os.RemoveAll(fpath)
		}
	}
	fmt.Println("\nCleanup complete!")
	fmt.Println("\ncd " + pluginName)
	fmt.Println("npm install")
	fmt.Println("composer install")
	fmt.Println("npm run dev")
}
