package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// define flags 
	helpFlag := flag.Bool("h", false, "see help messages")
	csvFlag := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limitFlag := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz order") 
	flag.Parse() 
	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0) 
	}
	
	// read in file 
	f, err := os.Open(*csvFlag)
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(f)
	checkMap := make(map[string]string)
	questions := make([]string, 0)
	for { 
		record, err := r.Read() 
		if err == io.EOF {
			break 
		} 
		if err != nil {
			panic(err) 
		}
		checkMap[record[0]] = record[1]
		questions = append(questions, record[0])
	}
	if *shuffle {
		rand.Seed(time.Now().Unix())
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	// start quiz 
	score := 0
	timer := time.NewTimer(time.Duration(*limitFlag) * time.Second)
	isTimeUp := false
	ansReader := bufio.NewReader(os.Stdin)  
	go func() {
		<- timer.C
		endRoutine(score, checkMap)
	}()
	for index, question := range questions {
		if isTimeUp {
			break 
		}
		fmt.Printf("Problem #%d: %v = ", index + 1, question)
		ans, err := ansReader.ReadString('\n') 
		if err != nil {
			panic(err)
		}
		ans = strings.TrimSpace(ans) 
		if ans == checkMap[question] {
			score += 1 
		}
	}
	endRoutine(score, checkMap)
}

func endRoutine(score int, checkMap map[string]string) {
	fmt.Printf("\nYou scored %d out of %d\n", score, len(checkMap))
	os.Exit(0)
}