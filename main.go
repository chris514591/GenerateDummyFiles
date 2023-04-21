package main

import (
	"io/fs"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	Lognconfig "GenerateDummyFiles/Packagelognconfig"

	Fileutil "GenerateDummyFiles/Packagefileutil"

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
	config, err := Fileutil.ReadConfigFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	errLogFile, err := Fileutil.OpenErrorsLogFile("errors.log")
	if err != nil {
		log.Fatalf("Failed to open errors.log file: %v", err)
	}
	defer errLogFile.Close()
	log.SetOutput(errLogFile)

	err = walkDirectories(config, generateFiles)
	if err != nil {
		log.Fatalf("Failed to walk directory: %v", err)
	}
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

func generateFiles(path string, numOfFiles int, fileExtensions []string, config Config) error {
	for i := 1; i <= numOfFiles; i++ {
		fileSize := rand.Intn(config.MaxFileSize-config.MinFileSize) + config.MinFileSize
		extension := fileExtensions[rand.Intn(len(fileExtensions))]
		fileName := Lognconfig.GenerateFileName(10)

		err := Fileutil.GenerateFile(filepath.Join(path, fileName+extension), lorem.Paragraph(1, fileSize/100))
		if err != nil {
			log.Printf("Failed to generate file: %v", err)
		}
	}

	return nil
}
