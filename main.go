package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
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
	scanner := bufio.NewScanner(os.Stdin)

	difficultySet := NewDifficultySet()
	difficulty, err := setupDifficulty(scanner, difficultySet)
	if err != nil {
		fmt.Printf("An error occurred during difficulty setup: %v\n", err)
		os.Exit(1)
	}

	secret := generateSecret(difficulty.RangeFrom, difficulty.RangeTo)
	attemp := 1
	var lastAttempts LastAttempts

	for {
		if len(lastAttempts) > 0 {
			color.Yellow("\nLast attempts: %s\n", lastAttempts.print())
		}

		color.Yellow("Try #%d - Enter a number from %d to %d:\n", attemp, difficulty.RangeFrom, difficulty.RangeTo)
		enteredNum, err := readGuess(scanner, difficulty.RangeFrom, difficulty.RangeTo)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("Game is over")
				return
			}
			fmt.Printf("An error occurred: %v\n", err)
			os.Exit(1)
		}

		direction, equal := compareGuess(enteredNum, secret)
		gameOver := equal || attemp == difficulty.Attempts
		if gameOver {
			outcome := "loss"
			if equal {
				color.Green("You win!")
				outcome = "win"
			} else {
				color.Red("Out of attempts. You lose... The secret number was: %d.", secret)
			}

			writeResult(outcome, attemp)

			attemp = 1
			lastAttempts = LastAttempts{}
			ok, err := startOver(scanner)
			if err != nil {
				fmt.Printf("An error occurred: %v\n", err)
				os.Exit(1)
			}
			if ok {
				difficulty, err = setupDifficulty(scanner, difficultySet)
				if err != nil {
					fmt.Printf("An error occurred during difficulty setup: %v\n", err)
					os.Exit(1)
				}
				secret = generateSecret(difficulty.RangeFrom, difficulty.RangeTo)
			} else {
				fmt.Println("Game is over")
				return
			}
		} else {
			attemp++
			lastAttempts = append(lastAttempts, enteredNum)
			fmt.Printf("%s The secret number is %s\n", proximityHint(enteredNum, secret), direction)
		}
	}
}

func setupDifficulty(scanner *bufio.Scanner, difficultySet DifficultySet) (Difficulty, error) {
	fmt.Println("Choose difficulty: 1 (easy), 2 (medium), 3 (hard)")

	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return Difficulty{}, fmt.Errorf("An error occurred: %w", err)
			}
			return Difficulty{}, fmt.Errorf("Unexpected end of input")
		}

		enteredDifficulty := scanner.Text()
		enteredDifficultyNum, err := strconv.Atoi(enteredDifficulty)
		allowedDifficultyNums := [3]int{1, 2, 3}
		if err != nil || !inArray(allowedDifficultyNums, enteredDifficultyNum) {
			fmt.Println("Please enter a valid number: 1 (easy), 2 (medium), 3 (hard)")
			continue
		}
		return difficultySet.getDifficulty(enteredDifficultyNum), nil
	}
}

func generateSecret(rangeFrom, rangeTo int) int {
	count := rangeTo - rangeFrom + 1
	return rand.Intn(count) + rangeFrom
}

func readGuess(scanner *bufio.Scanner, rangeFrom, rangeTo int) (int, error) {
	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return 0, err
			}
			return 0, io.EOF
		}

		enteredNum, err := strconv.Atoi(scanner.Text())
		if err != nil || enteredNum < rangeFrom || enteredNum > rangeTo {
			fmt.Printf("Please enter a number from %d to %d:\n", rangeFrom, rangeTo)
			continue
		}

		return enteredNum, nil
	}
}

func compareGuess(guess, secret int) (direction string, equal bool) {
	switch {
	case guess < secret:
		return "higher", false
	case guess > secret:
		return "lower", false
	default:
		return "", true
	}
}

func proximityHint(guess, secret int) string {
	distance := abs(guess - secret)

	switch {
	case distance <= 5:
		return "Hot!"
	case distance <= 15:
		return "Warm!"
	default:
		return "Cold."
	}
}

func startOver(scanner *bufio.Scanner) (bool, error) {
	fmt.Println("Want to play again? (y - yes, other button - no)")

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("An error occurred: %w", err)
		}
		return false, nil
	}

	if scanner.Text() == "y" {
		return true, nil
	}
	return false, nil
}

func writeResult(outcome string, attempts int) {
	const resultsFile = "results.json"

	results, err := readResults(resultsFile)
	if err != nil {
		fmt.Println("Error while saving result into file")
		return
	}

	results = append(results, Result{
		Date:     time.Now(),
		Outcome:  outcome,
		Attempts: attempts,
	})

	if err := saveResults(resultsFile, results); err != nil {
		fmt.Println("Error while saving result into file")
	}
}

func readResults(path string) ([]Result, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []Result{}, nil
	}
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []Result{}, nil
	}

	var results []Result
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func saveResults(path string, results []Result) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".results-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return err
	}

	return nil
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
