package quiz1

import "flag"

func RunQuiz1() {
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of question, answer")
	flag.Parse()
	_ = csvFileName
}