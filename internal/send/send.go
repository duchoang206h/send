package send

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/duchoang206h/send-cli/internal/api"
)

func validPath(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, os.ErrNotExist
	} else {
		return false, errors.New("error occurred while checking the path")
	}
}

func zipFolder(path, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	files := []string{}
	err = filepath.WalkDir(path, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, filePath)
		return nil
	})
	sort.Strings(files)
	for _, file := range files {
		relPath, err := filepath.Rel(path, file)
		if err != nil {
			return err
		}

		zipFileEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		_, err = io.Copy(zipFileEntry, fileToZip)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func zipFileName(path string) string {
	return fmt.Sprintf("%s.zip", filepath.Base(path))
}

func SendFile(ctx context.Context, path string) (string, error) {
	_, err := validPath(path)
	if err != nil {
		return "", err
	}
	response, err := api.UploadFile(ctx, path)
	fmt.Println("err", err)
	if err != nil {
		return "", err
	}
	return response.Result, nil
}

func SendDirectory(ctx context.Context, path string) (string, error) {
	_, err := validPath(path)
	if err != nil {
		return "", err
	}
	zipPath := zipFileName(path)
	err = zipFolder(path, zipPath)
	if err != nil {
		return "", err
	}
	result, err := SendFile(ctx, zipPath)
	if err != nil {
		return "", err
	}
	os.Remove(zipPath)
	return result, nil
}

func PrintResult(resChan <-chan string, errChan chan error) {
	animationChars := []string{"|", "/", "-", "\\"}
	for {
		select {
		case url := <-resChan:
			fmt.Print("\033[H\033[2J") // clear screen
			fmt.Println("File url: ", url)
			return
		case err := <-errChan:
			fmt.Print("\033[H\033[2J")
			fmt.Println("Ops!!!", err)
			return
		default:
			for _, char := range animationChars {
				fmt.Print("\033[H\033[2J")
				fmt.Print("Loading " + char)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
