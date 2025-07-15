package zip_service

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"github.com/google/uuid"
)

func DownloadAndArchive(urls []string, archivePath string, id uuid.UUID) error {


	tempDir, err := os.MkdirTemp("", "downloaded_files")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	var downloadedFiles []string
	var failed []string

	for _, url := range urls {

		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			failed = append(failed, fmt.Sprintf("%s (download failed)", url))
			continue
		}
		defer resp.Body.Close()

		fileName := filepath.Base(url)
		localPath := filepath.Join(tempDir, fileName)
		out, err := os.Create(localPath)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s (write failed)", url))
			continue
		}

		_, err = io.Copy(out, resp.Body)
		out.Close()
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s (copy failed)", url))
			continue
		}

		downloadedFiles = append(downloadedFiles, localPath)
	}

	if len(downloadedFiles) == 0 {
		return fmt.Errorf("all downloads failed: %v", failed)
	}

	err = createZip(downloadedFiles, archivePath, id)
	if err != nil {
		return fmt.Errorf("zip creation failed: %w", err)
	}

	if len(failed) > 0 {
		return fmt.Errorf("archive created with some failures: %v", failed)
	}

	return nil
}

func createZip(files []string, output string, id uuid.UUID) error {

	archivePath := filepath.Join(output, id.String()+".zip")
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, path := range files {
		fileToZip, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.Base(path)
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return err
		}
	}

	return nil
}

