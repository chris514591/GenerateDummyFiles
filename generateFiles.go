package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Path string `json:"path"`
}

func main() {
	minFileSize := 5000
	maxFileSize := 5000000

	fileExtensions := []string{".txt", ".cfg", ".csv"}

	rand.Seed(time.Now().UnixNano())
	numOfFiles := rand.Intn(51) + 50

	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	// Walk the given directory and its subdirectories
	err = filepath.Walk(config.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is a directory
		if info.IsDir() {
			// Generate files in this directory
			for i := 1; i <= numOfFiles; i++ {
				fileSize := rand.Intn(maxFileSize-minFileSize) + minFileSize

				extension := fileExtensions[rand.Intn(len(fileExtensions))]

				fileName := generateFileName()

				err := generateFile(path, fileName+extension, fileSize)
				if err != nil {
					panic(err)
				}
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func generateFileName() string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

	name := make([]byte, 10)
	for i := range name {
		name[i] = letters[rand.Intn(len(letters))]
	}

	return string(name)
}

func generateFile(path string, name string, size int) error {
	file, err := os.Create(filepath.Join(path, name))
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
