package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func extractTarGz(archivePath, outputPath string) error {
	// open archive
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	// find target file in archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Name == "system-link" {
			outputFile, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, tarReader)
			return err
		}
	}

	return fmt.Errorf("file 'system-link' not found in tar.gz archive")
}

func extractZip(archivePath, outputPath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.Name == "system-link.exe" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			outputFile, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, rc)
			return err
		}
	}

	return fmt.Errorf("file 'system-link.exe' not found in zip archive")
}
