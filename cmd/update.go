package cmd

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Koenigseder/badm/internal/core"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update BADM",
	Long:  `Update BADM to the newest version`,
	Run: func(_ *cobra.Command, _ []string) {
		err := update()
		if err != nil {
			fmt.Println("Error updating BADM:", err)
			os.Exit(1)
		}
	},
}

func update() error {
	filePermissions := os.FileMode(0755) //nolint:mnd

	// Get latest release information
	latestReleaseInformation, err := core.GetLatestReleaseInformation()
	if err != nil {
		return err
	}

	binaryDownloadURL := fmt.Sprintf(
		"https://github.com/Koenigseder/badm/releases/download/%s/badm_linux_amd64.tar.gz",
		latestReleaseInformation.TagName,
	)

	// Download archive
	tempArchivePath := fmt.Sprintf("badm_%s.tar.gz", latestReleaseInformation.TagName)
	if err = downloadFile(binaryDownloadURL, tempArchivePath); err != nil {
		return fmt.Errorf("error downloading binary archive: %v", err)
	}

	// Extract archive
	tempDir := fmt.Sprintf("badm_%s", latestReleaseInformation.TagName)
	if err = os.MkdirAll(tempDir, filePermissions); err != nil {
		return fmt.Errorf("error creating temp directory: %v", err)
	}

	if err = extractTarGz(tempArchivePath, tempDir); err != nil {
		return fmt.Errorf("error extracting archive: %v", err)
	}

	// Replace old binary with new one
	// Linux: os.Executable is the path of the running binary
	oldPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting currently running binary: %v", err)
	}

	newBinaryPath := filepath.Join(tempDir, "badm")
	if err = os.Rename(newBinaryPath, oldPath); err != nil {
		return fmt.Errorf("error replacing binary: %v", err)
	}

	// Set permissions
	if err = os.Chmod(oldPath, filePermissions); err != nil {
		return fmt.Errorf("error setting binary permissions: %v", err)
	}

	// Remove resources
	if err = os.RemoveAll(tempDir); err != nil {
		return fmt.Errorf("error removing temp directory: %v", err)
	}

	if err = os.RemoveAll(tempArchivePath); err != nil {
		return fmt.Errorf("error removing binary tarball: %v", err)
	}

	return nil
}

func downloadFile(url, outputPath string) error {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}

func extractTarGz(archivePath, targetDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}

	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}

	defer gz.Close()

	tarReader := tar.NewReader(gz)

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, header.Name) //nolint:gosec

		filePermissions := os.FileMode(0755) //nolint:mnd

		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(targetPath, filePermissions); err != nil {
				return err
			}

		case tar.TypeReg:
			if err = os.MkdirAll(filepath.Dir(targetPath), filePermissions); err != nil {
				return err
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				return err
			}

			defer outFile.Close()

			if _, err = io.Copy(outFile, tarReader); err != nil { //nolint:gosec
				return err
			}
		}
	}

	return nil
}
