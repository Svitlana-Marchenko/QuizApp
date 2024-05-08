package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	fName := flag.String("f", "quiz.csv", "file name to read quiz from")
	timer := flag.Int("t", 300, "timer in seconds for the quiz")
	flag.Parse()

	problems, err := problemPuller(*fName)
	if err != nil {
		exit(fmt.Sprintf("Smth went wrong: %s", err.Error()))
	}

	correctAnsw := 0

	timerObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansChan := make(chan string)

problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("%d: %s=", i+1, p.question)

		go func() {
			fmt.Scanln(&answer)
			ansChan <- answer
		}()
		select {
		case <-timerObj.C:
			fmt.Println("Timer expired")
			break problemLoop
		case isAnswer := <-ansChan:
			if isAnswer == p.answer {
				correctAnsw++
			}
			if i == len(problems)-1 {
				close(ansChan)
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correctAnsw, len(problems))
	fmt.Println("Press enter to quit")
}

func problemPuller(filename string) ([]problem, error) {

	if file, err := os.Open(filename); err == nil {
		csvReader := csv.NewReader(file)
		if all, err := csvReader.ReadAll(); err == nil {
			return problemParser(all), nil
		} else {
			return nil, fmt.Errorf("error while parsing %s: %s", filename, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error while opening %s: %s", filename, err.Error())
	}
}

func problemParser(lines [][]string) []problem {
	problems := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		problems[i] = problem{lines[i][0], lines[i][1]}
	}
	return problems
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

type problem struct {
	question string
	answer   string
}
