package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
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

	fmt.Printf("Starting Quiz Game with args: %s File with a time limit of %d Seconds\n", csvFlag, limitFlag)

	game_records := readCsvFile(csvFlag)
	correct := 0

	timer := time.NewTimer(time.Duration(limitFlag) * time.Second)

	for _, record := range game_records {
		fmt.Print(fmt.Sprintf("%s: ", record[0]))
		answerChan := make(chan string)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			answerChan <- strings.TrimSpace(text)
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nTotal Score: %d out of %d Correct\n", correct, len(game_records))
			return
		case answer := <-answerChan:
			if answer == record[1] {
				correct += 1
			}
		}
	}

	fmt.Printf("Total Score: %d out of %d Correct\n", correct, len(game_records))
}
