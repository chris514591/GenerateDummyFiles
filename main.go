package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	lorem "github.com/drhodes/golorem"
)

type Config struct {
	Path          string `json:"path"`
	MinNumOfFiles int    `json:"min_num_of_files"`
	MaxNumOfFiles int    `json:"max_num_of_files"`
	MinFileSize   int    `json:"min_file_size"`
	MaxFileSize   int    `json:"max_file_size"`
}

func main() {
	config, err := readConfigFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	errLogFile, err := openErrorsLogFile("errors.log")
	if err != nil {
		log.Fatalf("Failed to open errors.log file: %v", err)
	}
	defer errLogFile.Close()
	log.SetOutput(errLogFile)

	var totalNumOfFiles, numOfGeneratedFilesTotal int
	err = walkDirectories(config, func(path string, numOfFiles int, fileExtensions []string, config Config) error {
		totalNumOfFiles += numOfFiles
		return generateFiles(path, numOfFiles, fileExtensions, config, &numOfGeneratedFilesTotal, totalNumOfFiles)
	})
	if err != nil {
		log.Fatalf("Failed to walk directory: %v", err)
	}

	fmt.Printf("Generated all files (%d/%d, %.0f%%)\n", numOfGeneratedFilesTotal, totalNumOfFiles, float64(numOfGeneratedFilesTotal)/float64(totalNumOfFiles)*100)
}

func readConfigFile(filename string) (Config, error) {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func openErrorsLogFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}

func walkDirectories(config Config, generateFilesFunc func(string, int, []string, Config) error) error {
	rand.Seed(time.Now().UnixNano())

	return filepath.WalkDir(config.Path, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Failed to walk directory: %v", err)
			return nil
		}

		if dirEntry.IsDir() {
			numOfFiles := rand.Intn(config.MaxNumOfFiles-config.MinNumOfFiles+1) + config.MinNumOfFiles
			fileExtensions := []string{".txt", ".csv", ".html"}

			return generateFilesFunc(path, numOfFiles, fileExtensions, config)
		}

		return nil
	})
}

func generateFiles(path string, numOfFiles int, fileExtensions []string, config Config, numOfGeneratedFilesTotal *int, totalNumOfFiles int) error {
	for i := 1; i <= numOfFiles; i++ {
		fileSize := rand.Intn(config.MaxFileSize-config.MinFileSize) + config.MinFileSize
		extension := fileExtensions[rand.Intn(len(fileExtensions))]
		fileName := generateFileName(10)

		err := generateFile(filepath.Join(path, fileName+extension), lorem.Paragraph(1, fileSize/100))
		if err != nil {
			log.Printf("Failed to generate file: %v", err)
		} else {
			*numOfGeneratedFilesTotal++
			percentage := float64(*numOfGeneratedFilesTotal) / float64(totalNumOfFiles) * 100
			fmt.Printf("Generated %d/%d files (%.0f%%)\n", *numOfGeneratedFilesTotal, totalNumOfFiles, percentage)
		}
	}

	return nil
}

func generateFileName(length int) string {
	const characters = "abcdefghijklmnopqrstuvwxyz0123456789"

	fileName := make([]byte, length)
	for i := range fileName {
		fileName[i] = characters[rand.Intn(len(characters))]
	}

	return string(fileName)
}

func generateFile(filePath string, data string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
