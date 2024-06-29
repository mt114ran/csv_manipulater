package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getCSVFiles() ([]string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var csvFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			csvFiles = append(csvFiles, file.Name())
		}
	}

	return csvFiles, nil
}

func splitCSV(filePath string, linesPerFile int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")

	for part := 1; ; part++ {
		outputFileName := fmt.Sprintf("split_%s_%d.csv", timestamp, part)
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			return err
		}
		writer := csv.NewWriter(outputFile)
		writer.Write(headers)

		linesWritten := 0
		for i := 0; i < linesPerFile; i++ {
			record, err := reader.Read()
			if err != nil {
				if err.Error() == "EOF" {
					writer.Flush()
					outputFile.Close()
					if linesWritten == 0 {
						os.Remove(outputFileName)
					}
					return nil
				}
				return err
			}
			writer.Write(record)
			linesWritten++
		}
		writer.Flush()
		outputFile.Close()
	}
}

func extractLines(filePath string, lines int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	outputFileName := fmt.Sprintf("extracted_%s.csv", timestamp)
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	writer.Write(headers)

	for i := 0; i < lines; i++ {
		record, err := reader.Read()
		if err != nil {
			break
		}
		writer.Write(record)
	}

	writer.Flush()
	return nil
}

func main() {
	fmt.Print("Select a function to execute (1: Split, 2: Extract): ")
	var functionChoice int
	fmt.Scan(&functionChoice)

	csvFiles, err := getCSVFiles()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(csvFiles) == 0 {
		fmt.Println("No CSV files found.")
		return
	}

	fmt.Println("The following CSV files were found:")
	for i, file := range csvFiles {
		fmt.Printf("%d: %s\n", i+1, file)
	}

	fmt.Print("Enter the number of the CSV file to process: ")
	var choice int
	fmt.Scan(&choice)

	if choice < 1 || choice > len(csvFiles) {
		fmt.Println("Invalid selection.")
		return
	}

	selectedFile := csvFiles[choice-1]

	reader := bufio.NewReader(os.Stdin)

	if functionChoice == 1 {
		fmt.Print("Enter the number of lines per split: ")
		lineInput, _ := reader.ReadString('\n')
		linesPerFile, err := strconv.Atoi(strings.TrimSpace(lineInput))
		if err != nil {
			log.Fatalf("Invalid input: %v", err)
		}

		err = splitCSV(selectedFile, linesPerFile)
		if err != nil {
			log.Fatalf("Error: %v", err)
		} else {
			fmt.Println("Splitting completed.")
		}
	} else if functionChoice == 2 {
		fmt.Print("Enter the number of lines to extract: ")
		lineInput, _ := reader.ReadString('\n')
		lines, err := strconv.Atoi(strings.TrimSpace(lineInput))
		if err != nil {
			log.Fatalf("Invalid input: %v", err)
		}

		err = extractLines(selectedFile, lines)
		if err != nil {
			log.Fatalf("Error: %v", err)
		} else {
			fmt.Println("Extraction completed.")
		}
	} else {
		fmt.Println("Invalid selection.")
	}
}
