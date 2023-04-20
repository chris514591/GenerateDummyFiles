package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
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

	for i := 1; i <= numOfFiles; i++ {
		fileSize := rand.Intn(maxFileSize-minFileSize) + minFileSize

		extension := fileExtensions[rand.Intn(len(fileExtensions))]

		fileName := "testfile" + strconv.Itoa(i) + extension

		err := generateFile(config.Path, fileName, fileSize)
		if err != nil {
			panic(err)
		}
	}
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
