package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

const (
	defaultRangeFrom int = 1
	defaultRangeTo   int = 100
	defaultAttemps   int = 10
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

func (d DifficultySet) get(difficulty int) Difficulty {
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

func main() {
	difficultySet := NewDifficultySet()
	difficulty := Difficulty{
		RangeFrom: defaultRangeFrom,
		RangeTo:   defaultRangeTo,
		Attempts:  defaultAttemps,
	}

	randNum := generateRandNum(defaultRangeFrom, defaultRangeTo)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Let's start the game! Choose difficulty: 1 (easy), 2 (medium), 3 (hard)")

	for {
		if scanner.Scan() {
			enteredDifficulty := scanner.Text()
			enteredDifficultyNum, err := strconv.Atoi(enteredDifficulty)
			allowedDifficultyNums := [3]int{1, 2, 3}
			if err != nil || !inArray(allowedDifficultyNums, enteredDifficultyNum) {
				fmt.Println("Please enter a valid number: 1 (easy), 2 (medium), 3 (hard)")
				continue
			}
			difficulty = difficultySet.get(enteredDifficultyNum)
			randNum = generateRandNum(difficulty.RangeFrom, difficulty.RangeTo)
			break
		}
	}

	for {
		fmt.Printf("Enter a number from %d to %d!\n", difficulty.RangeFrom, difficulty.RangeTo)
		if scanner.Scan() {
			enteredString := scanner.Text()
			enteredNum, err := strconv.Atoi(enteredString)
			if err != nil {
				fmt.Printf("Please enter a number from %d to %d!\n", difficulty.RangeFrom, difficulty.RangeTo)
				continue
			}

			if enteredNum == randNum {
				fmt.Println("You win! Want to play again? (y - yes, n - no)")
				if scanner.Scan() {
					if scanner.Text() == "y" {
						randNum = generateRandNum(defaultRangeFrom, defaultRangeTo) // TODO
						continue
					} else if scanner.Text() == "n" {
						fmt.Println("Game is over")
						return
					}
				}
			} else {
				fmt.Println(giveHint(enteredNum, randNum))
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("An error occurred: %v\n", err)
		}
	}
}

func generateRandNum(rangeFrom, rangeTo int) int {
	count := rangeTo - rangeFrom + 1
	return rand.Intn(count) + rangeFrom
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
