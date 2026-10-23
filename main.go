package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type DifficultySet struct {
	Easy   Difficulty
	Medium Difficulty
	Hard   Difficulty
}

type Difficulty struct {
	RangeFrom int
	RangeTo   int
	Attempts  int
}

func NewDifficultySet() DifficultySet {
	return DifficultySet{
		Easy:   Difficulty{RangeFrom: 1, RangeTo: 50, Attempts: 15},
		Medium: Difficulty{RangeFrom: 1, RangeTo: 100, Attempts: 10},
		Hard:   Difficulty{RangeFrom: 1, RangeTo: 200, Attempts: 5},
	}
}

func (d DifficultySet) getDifficulty(difficulty int) Difficulty {
	difficultyMap := map[int]Difficulty{
		1: d.Easy,
		2: d.Medium,
		3: d.Hard,
	}

	val, ok := difficultyMap[difficulty]
	if !ok {
		return d.Medium
	}

	return val
}

type Result struct {
	Date     time.Time `json:"date"`
	Outcome  string    `json:"outcome"`
	Attempts int       `json:"attempts"`
}

type LastAttempts []int

func (la LastAttempts) print() string {
	var sb strings.Builder

	sb.WriteString(strconv.Itoa(la[0]))

	for _, attempt := range la[1:] {
		sb.WriteString(",")
		sb.WriteString(strconv.Itoa(attempt))
	}

	return sb.String()
}

func main() {
	difficultySet := NewDifficultySet()
	difficulty, err := setupDifficulty(difficultySet)
	if err != nil {
		fmt.Printf("An error occurred during difficulty setup: %v\n", err)
		os.Exit(1)
	}

	randNum := generateRandNum(difficulty.RangeFrom, difficulty.RangeTo)
	attemp := 1
	var lastAttempts LastAttempts

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if len(lastAttempts) > 0 {
			fmt.Printf("\nLast attempts: %s\n", lastAttempts.print())
		}

		fmt.Printf("Try #%d - Enter a number from %d to %d:\n", attemp, difficulty.RangeFrom, difficulty.RangeTo)
		if scanner.Scan() {
			enteredString := scanner.Text()
			enteredNum, err := strconv.Atoi(enteredString)
			outOfRange := enteredNum < difficulty.RangeFrom || enteredNum > difficulty.RangeTo
			if err != nil || outOfRange {
				fmt.Printf("Please enter a number from %d to %d:\n", difficulty.RangeFrom, difficulty.RangeTo)
				continue
			}

			gameOver := enteredNum == randNum || attemp == difficulty.Attempts
			if gameOver {
				outcome := "loss"
				if enteredNum == randNum {
					fmt.Println("You win!")
					outcome = "win"
				} else {
					fmt.Println("Out of attempts. You lose...")
				}

				writeResult(outcome, attemp)

				attemp = 1
				lastAttempts = LastAttempts{}
				ok, err := startOver()
				if err != nil {
					fmt.Printf("An error occurred: %v\n", err)
					os.Exit(1)
				}
				if ok {
					difficulty, err = setupDifficulty(difficultySet)
					if err != nil {
						fmt.Printf("An error occurred during difficulty setup: %v\n", err)
						os.Exit(1)
					}
					randNum = generateRandNum(difficulty.RangeFrom, difficulty.RangeTo)
				} else {
					fmt.Println("Game is over")
					return
				}
			} else {
				attemp++
				lastAttempts = append(lastAttempts, enteredNum)
				fmt.Println(giveHint(enteredNum, randNum))
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("An error occurred: %v\n", err)
			os.Exit(1)
		}
	}
}

func setupDifficulty(difficultySet DifficultySet) (Difficulty, error) {
	fmt.Println("Choose difficulty: 1 (easy), 2 (medium), 3 (hard)")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		if scanner.Scan() {
			enteredDifficulty := scanner.Text()
			enteredDifficultyNum, err := strconv.Atoi(enteredDifficulty)
			allowedDifficultyNums := [3]int{1, 2, 3}
			if err != nil || !inArray(allowedDifficultyNums, enteredDifficultyNum) {
				fmt.Println("Please enter a valid number: 1 (easy), 2 (medium), 3 (hard)")
				continue
			}
			return difficultySet.getDifficulty(enteredDifficultyNum), nil
		}

		if err := scanner.Err(); err != nil {
			return Difficulty{}, fmt.Errorf("An error occurred: %w", err)
		}
	}
}

func generateRandNum(rangeFrom, rangeTo int) int {
	count := rangeTo - rangeFrom + 1
	return rand.Intn(count) + rangeFrom
}

func startOver() (bool, error) {
	fmt.Println("Want to play again? (y - yes, other button - no)")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		if scanner.Text() == "y" {
			return true, nil
		}
		return false, nil
	}
	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("An error occurred: %w", err)
	}

	return false, nil
}

func giveHint(enteredNum, randNum int) string {
	hint := "higher"
	if enteredNum > randNum {
		hint = "lower"
	}

	proximity := "Cold."
	distance := abs(enteredNum - randNum)

	switch {
	case distance <= 5:
		proximity = "Hot!"
	case distance <= 15:
		proximity = "Warm!"
	}

	return fmt.Sprintf("%s The secret number is %s", proximity, hint)
}

func writeResult(outcome string, attempts int) {
	file, err := os.OpenFile("results.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Print("Error while saving result into file")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(Result{
		Date:     time.Now(),
		Outcome:  outcome,
		Attempts: attempts,
	})
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func inArray(arr [3]int, target int) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}
