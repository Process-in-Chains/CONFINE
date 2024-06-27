package reset

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ReplaceWithEmptyMap(filePath string) error {
	// Create an empty map as a JSON string
	emptyMapJSON := "{}"

	// Write the empty map JSON to the specified file
	err := ioutil.WriteFile(filePath, []byte(emptyMapJSON), 0644)
	if err != nil {
		return err
	}
	return nil
}
func ReplaceWithEmptyMatrix(filePath string) error {
	// Create an empty map as a JSON string
	emptyMapJSON := "[[]]"

	// Write the empty map JSON to the specified file
	err := ioutil.WriteFile(filePath, []byte(emptyMapJSON), 0644)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAllFilesInSubfolders(rootDir string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {

			return nil
		}
		if !info.IsDir() {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func DeleteEmptySubfolders(rootDir string) error {
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {

			return nil
		}
		if info.IsDir() {
			// Check if the directory is empty
			isEmpty, _ := isDirectoryEmpty(path)
			if isEmpty {

				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// Helper function to check if a directory is empty
func isDirectoryEmpty(path string) (bool, error) {
	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(dirEntries) == 0, nil
}
func DeleteTraceFolders() error {
	rootDir := "mining-data/provision-data/process-01/"
	// Get a list of all files and folders in the root directory
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return err
	}
	// Iterate over each file/folder
	for _, file := range files {
		// Check if it's a directory and its name starts with "trace_"
		if file.IsDir() && strings.HasPrefix(file.Name(), "trace_") {
			DeleteAllFilesInSubfolders(filepath.Join(rootDir, file.Name()))
			// Get the full path of the directory
			dirPath := filepath.Join(rootDir, file.Name())
			// Delete the directory and all its contents
			err := os.Remove(dirPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
