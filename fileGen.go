package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Config represents the configuration data read from the JSON file.
type Config struct {
	// Comments aren't allowed in JSON files so: Use forward slashes, not back slashes for the folder path!
	Path           string   `json:"path"`
	MinNumOfFiles  int      `json:"min_num_of_files"`
	MaxNumOfFiles  int      `json:"max_num_of_files"`
	MinFileSize    int      `json:"min_file_size"`
	MaxFileSize    int      `json:"max_file_size"`
	FileExtensions []string `json:"file_extensions"`
}

func main() {
	// Read the configuration from the config file
	config, err := readConfigFile("config.json")
	if err != nil {
		log.Printf("Failed to read config file: %v", err)
	}

	// Open the errors.log file for logging errors
	errLogFile, err := openErrorsLogFile("errors.log")
	if err != nil {
		fmt.Printf("Failed to open errors.log file: %v", err)
	}
	defer errLogFile.Close()
	log.SetOutput(errLogFile)

	var totalNumOfFiles, numOfGeneratedFilesTotal int
	var totalFilesSize int64

	// Measure start time
	startTime := time.Now()

	// Walk through directories and generate files
	err = walkDirectories(config, func(path string, numOfFiles int, config Config) error {
		totalNumOfFiles += numOfFiles
		return generateFiles(path, numOfFiles, config, &numOfGeneratedFilesTotal, &totalFilesSize, totalNumOfFiles, startTime)
	})
	if err != nil {
		log.Printf("Failed to walk directory: %v", err)
	}

	// Measure end time
	endTime := time.Now()

	// Calculate and print elapsed time in seconds
	elapsedTime := endTime.Sub(startTime).Seconds()
	fmt.Printf("Elapsed time: %.2f seconds\n", elapsedTime)

	// Print the summary of generated files
	fmt.Printf("Generated all files (%d/%d, %.2f%%)\n", numOfGeneratedFilesTotal, totalNumOfFiles, float64(numOfGeneratedFilesTotal)/float64(totalNumOfFiles)*100)
	printFileSize("Total size of generated files", totalFilesSize)

	// Wait for user input before exiting
	waitForEnterKey()
}

// waitForEnterKey waits for the user to press Enter before exiting the program.
func waitForEnterKey() {
	fmt.Println("Press Enter to exit...")
	_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Printf("Failed to read from stdin: %v", err)
	}
}

// readConfigFile reads the configuration data from the provided JSON file.
func readConfigFile(filename string) (Config, error) {
	configFile, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Failed to read config file: %v", err)
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Printf("Failed to unmarshal config file: %v", err)
	}

	return config, nil
}

// openErrorsLogFile opens the errors.log file for writing error logs.
func openErrorsLogFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open errors.log file: %v", err)
	}

	return file, nil
}

// walkDirectories recursively walks through directories starting from the provided path.
func walkDirectories(config Config, generateFilesFunc func(string, int, Config) error) error {
	rand.Seed(time.Now().UnixNano())

	return filepath.WalkDir(config.Path, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Failed to walk directory: %v", err)
			return nil
		}

		if dirEntry.IsDir() {
			// Generate random number of files
			numOfFiles := rand.Intn(config.MaxNumOfFiles-config.MinNumOfFiles+1) + config.MinNumOfFiles

			// Generate files in the current directory
			return generateFilesFunc(path, numOfFiles, config)
		}

		return nil
	})
}

// generateFiles generates the specified number of files in the given path with random content.
func generateFiles(path string, numOfFiles int, config Config, numOfGeneratedFilesTotal *int, totalFilesSize *int64, totalNumOfFiles int, startTime time.Time) error {
	for i := 1; i <= numOfFiles; i++ {
		// Generate random file size, extension, and file name
		fileSize := rand.Intn(config.MaxFileSize-config.MinFileSize) + config.MinFileSize
		extension := config.FileExtensions[rand.Intn(len(config.FileExtensions))]
		fileName := generateFileName(10)

		// Generate a file with random content
		err := generateFile(filepath.Join(path, fileName+extension), fileSize)
		if err != nil {
			log.Printf("Failed to generate file: %v", err)
		} else {
			*numOfGeneratedFilesTotal++
			*totalFilesSize += int64(fileSize)
			percentage := float64(*numOfGeneratedFilesTotal) / float64(totalNumOfFiles) * 100
			elapsedTime := time.Since(startTime).Seconds()
			fmt.Printf("Generated %d/%d files (%.2f%%) | Elapsed time: %.2f seconds\n", *numOfGeneratedFilesTotal, totalNumOfFiles, percentage, elapsedTime)
		}
	}

	return nil
}

// generateFileName generates a random file name with the specified length.
func generateFileName(length int) string {
	const characters = "abcdefghijklmnopqrstuvwxyz0123456789"

	fileName := make([]byte, length)
	for i := range fileName {
		fileName[i] = characters[rand.Intn(len(characters))]
	}

	return string(fileName)
}

// generateFile creates a file with the specified file path and writes the provided data into it.
func generateFile(filePath string, fileSize int) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	// Write "Success" at the beginning of the file
	_, err = file.WriteString("Success")
	if err != nil {
		log.Printf("Failed to write data to file: %v", err)
		return err
	}

	// Calculate the remaining size after writing "Success"
	remainingSize := fileSize - len("Success")

	// Create a byte slice of the remaining size filled with zeroes
	data := make([]byte, remainingSize)

	// Write the data to the file
	_, err = file.Write(data)
	if err != nil {
		log.Printf("Failed to write data to file: %v", err)
		return err
	}

	return nil
}

// Function to print file size in bytes, megabytes, and gigabytes
func printFileSize(label string, size int64) {
	fmt.Printf("%s: %d bytes\n", label, size)
	fmt.Printf("%s: %.2f MB\n", label, float64(size)/(1024*1024))
	fmt.Printf("%s: %.2f GB\n", label, float64(size)/(1024*1024*1024))
}
