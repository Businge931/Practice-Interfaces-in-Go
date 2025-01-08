package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileSystemDatabase struct {
	BaseDir string
}

func NewFileSystemDatabase(baseDir string) *FileSystemDatabase {
	return &FileSystemDatabase{BaseDir: baseDir}
}

func (fs *FileSystemDatabase) Create(location string, data map[string]interface{}) (bool, string) {
	filePath := filepath.Join(fs.BaseDir, location)

	// Check if the file already exists.
	if _, err := os.Stat(filePath); err == nil {
		return false, "File already exists"
	}

	// Ensure the directory structure exists.
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return false, fmt.Sprintf("Failed to create directories: %v", err)
	}

	// Serialize the data to JSON.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return false, fmt.Sprintf("Failed to serialize data: %v", err)
	}

	// Write the data to the file.
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return false, fmt.Sprintf("Failed to write file: %v", err)
	}

	return true, ""
}

func (fs *FileSystemDatabase) Read(location string) (bool, string, map[string]interface{}) {
	filePath := filepath.Join(fs.BaseDir, location)

	// Check if the file exists.
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false, "File does not exist", nil
	}

	// Read the file contents.
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Sprintf("Failed to read file: %v", err), nil
	}

	// Deserialize the JSON data.
	var data map[string]interface{}
	if err := json.Unmarshal(fileData, &data); err != nil {
		return false, fmt.Sprintf("Failed to parse JSON: %v", err), nil
	}

	return true, "", data
}

func (fs *FileSystemDatabase) Update(location string, data map[string]interface{}) (bool, string) {
	filePath := filepath.Join(fs.BaseDir, location)

	// Check if the file exists.
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false, "File does not exist"
	}

	// Serialize the data to JSON.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return false, fmt.Sprintf("Failed to serialize data: %v", err)
	}

	// Write the data to the file.
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return false, fmt.Sprintf("Failed to write file: %v", err)
	}

	return true, ""
}

func (fs *FileSystemDatabase) Delete(location string) (bool, string) {
	filePath := filepath.Join(fs.BaseDir, location)

	// Check if the file exists.
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false, "File does not exist"
	}

	// Delete the file.
	if err := os.Remove(filePath); err != nil {
		return false, fmt.Sprintf("Failed to delete file: %v", err)
	}

	return true, ""
}
