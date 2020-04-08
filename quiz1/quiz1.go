package quiz1

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func RunQuiz1() {
	csvFileName := flag.String("csv", "quiz1/problems.csv", "a csv file in the format of question, answer")
	timeLimit := flag.Int("limit", 2, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFileName))
	}
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i + 1, p.question)
		answerCh := make(chan string)
		go func() {
			var answer string
			_, err := fmt.Scanf("%s\n", &answer)
			if err != nil {
				fmt.Println("Failed to scan input!")
			}
			answerCh <- answer
		}()
		select {
		case <- timer.C:
			fmt.Println()
			break problemLoop
		case answer := <- answerCh:
			if answer == p.answer {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{question:strings.TrimSpace(line[0]), answer:strings.TrimSpace(line[1])}
	}
	return ret
}

type problem struct {
	question string
	answer string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}