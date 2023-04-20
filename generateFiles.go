package main

import (
	"encoding/json"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Path          string `json:"path"`
	MinNumOfFiles int    `json:"min_num_of_files"`
	MaxNumOfFiles int    `json:"max_num_of_files"`
}

func main() {
	minFileSize := 5000
	maxFileSize := 5000000

	fileExtensions := []string{".txt", ".cfg", ".csv"}

	rand.Seed(time.Now().UnixNano())

	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}

	numOfFiles := rand.Intn(config.MaxNumOfFiles-config.MinNumOfFiles+1) + config.MinNumOfFiles

	// Check if errors.log file exists, and create it if it doesn't
	if _, err := os.Stat("errors.log"); os.IsNotExist(err) {
		if _, err := os.Create("errors.log"); err != nil {
			log.Fatalf("Failed to create errors.log file: %v", err)
		}
	}

	err = filepath.WalkDir(config.Path, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dirEntry.IsDir() {
			for i := 1; i <= numOfFiles; i++ {
				fileSize := rand.Intn(maxFileSize-minFileSize) + minFileSize

				extension := fileExtensions[rand.Intn(len(fileExtensions))]

				fileName := generateFileName()

				err := generateFile(path, fileName+extension, fileSize)
				if err != nil {
					// Log the error to errors.log file
					logFile, err := os.OpenFile("errors.log", os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						log.Fatalf("Failed to open errors.log file: %v", err)
					}
					defer logFile.Close()
					log.SetOutput(logFile)
					log.Printf("Failed to generate file: %v", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Failed to walk directory: %v", err)
	}
}

func generateFileName() string {
	const characters = "abcdefghijklmnopqrstuvwxyz0123456789"

	fileName := make([]byte, 10)
	for i := range fileName {
		fileName[i] = characters[rand.Intn(len(characters))]
	}

	return string(fileName)
}

func generateFile(path string, fileName string, size int) error {
	file, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]byte, size)
	rand.Read(data)

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
