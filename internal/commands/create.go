package commands

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/prappo/plugkit/internal/config"
	"github.com/schollz/progressbar/v3"
)

func CreatePlugin(pluginName string) error {
	fmt.Printf("Creating WordPress plugin: %s\n", pluginName)

	// Download boilerplate
	zipPath := pluginName + ".zip"
	if err := downloadBoilerplate(zipPath); err != nil {
		return fmt.Errorf("failed to download boilerplate: %w", err)
	}
	defer os.Remove(zipPath)

	// Extract files
	if err := extractFiles(zipPath, pluginName); err != nil {
		return fmt.Errorf("failed to extract files: %w", err)
	}

	fmt.Printf("\nPlugin '%s' created successfully! 🎉\n", pluginName)
	cleanUp(pluginName)
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
}
