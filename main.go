package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	fileName := flag.String("f", "quiz.csv", "file name to read quiz from")
	timer := flag.Int("t", 300, "timer in seconds for the quiz")
	flag.Parse()

	problems, err := problemPuller(*fileName)
	if err != nil {
		exit(fmt.Sprintf("Smth went wrong: %s", err.Error()))
	}

	var quantity int
	for {
		fmt.Println("Enter the number of questions you want to quiz:")
		_, err = fmt.Scanln(&quantity)
		if err != nil {
			exit(fmt.Sprintf("Smth went wrong: %s", err.Error()))
		}
		if quantity <= len(problems) && quantity > 0 {
			break
		}
	}

	problemList := getRandomProblems(quantity, problems)

	correctAnsw := 0

	timerObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansChan := make(chan string)

problemLoop:
	for i, p := range problemList {
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
			if i == len(problemList)-1 {
				close(ansChan)
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correctAnsw, len(problemList))
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

func getRandomProblems(quantity int, problems []problem) []problem {
	if quantity >= len(problems) {
		return problems
	}

	indexes := make([]int, quantity)
	for i := 0; i < quantity; i++ {
		indexes[i] = rand.Intn(len(problems))
	}

	rndProblems := make([]problem, quantity)
	for i := 0; i < quantity; i++ {
		rndProblems[i] = problems[indexes[i]]
	}

	return rndProblems
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
