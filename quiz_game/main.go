package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var helpFlag bool
var limitFlag int
var csvFlag string

func readCsvFile(filePath string) [][]string {
	quiz_file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Fatal error opening file: "+filePath, err)
	}
	defer quiz_file.Close()

	csvReader := csv.NewReader(quiz_file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {

	flag.IntVar(&limitFlag, "limit", 30, "Integer value for the quiz time limit")
	flag.StringVar(&csvFlag, "csv", "problems.csv", "A CSV Filepath")

	flag.Parse()

	fmt.Println(fmt.Sprintf("Starting Quiz Game with args: %s File with a time limit of %d Seconds", csvFlag, limitFlag))

	game_records := readCsvFile(csvFlag)
	correct := 0
	for _, record := range game_records {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(fmt.Sprintf("%s: ", record[0]))
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == record[1] {
			correct += 1
		}
	}

	fmt.Println(fmt.Sprintf("Total Score: %d out of %d Correct", correct, len(game_records)))
}
