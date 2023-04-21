package main

import (
	"encoding/json"
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
	fileExtensions := []string{".txt", ".cfg", ".csv"}

	rand.Seed(time.Now().UnixNano())

	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Printf("Failed to read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Printf("Failed to unmarshal config file: %v", err)
	}

	numOfFiles := rand.Intn(config.MaxNumOfFiles-config.MinNumOfFiles+1) + config.MinNumOfFiles

	logFile, err := os.OpenFile("errors.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to open errors.log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	err = filepath.WalkDir(config.Path, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Failed to walk directory: %v", err)
			return nil
		}

		if dirEntry.IsDir() {
			for i := 1; i <= numOfFiles; i++ {
				fileSize := rand.Intn(config.MaxFileSize-config.MinFileSize) + config.MinFileSize

				extension := fileExtensions[rand.Intn(len(fileExtensions))]

				fileName := generateFileName()

				err := generateFile(path, fileName+extension, fileSize)
				if err != nil {
					log.Printf("Failed to generate file: %v", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("Failed to walk directory: %v", err)
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

	data := []byte(lorem.Paragraph(1, size/100))

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
