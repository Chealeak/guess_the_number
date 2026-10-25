# Guess the Number

A command-line game written in Go. The program picks a secret number, and you try to find it before you run out of attempts. After each wrong guess it tells you whether the secret number is higher or lower, and how close you are.

## How to play

1. Choose a difficulty:
   - **1 — Easy:** numbers from 1 to 50, 15 attempts
   - **2 — Medium:** numbers from 1 to 100, 10 attempts
   - **3 — Hard:** numbers from 1 to 200, 5 attempts
2. Enter a guess in the allowed range.
3. Read the hint:
   - **Hot!** — within 5 of the secret number
   - **Warm!** — within 15
   - **Cold.** — farther away
4. Previous guesses are listed before each new try.
5. When the round ends, the result is appended to `results.json`.
6. Type `y` to play again, or any other key to quit.

## Requirements

- [Go](https://go.dev/dl/) 1.26 or later

## How to run

From the project directory:

```bash
go run main.go
```

## Sample output

```text
Choose difficulty: 1 (easy), 2 (medium), 3 (hard)
2
Try #1 - Enter a number from 1 to 100:
50
Warm! The secret number is higher

Last attempts: 50
Try #2 - Enter a number from 1 to 100:
70
Hot! The secret number is lower

Last attempts: 50,70
Try #3 - Enter a number from 1 to 100:
65
You win!
Want to play again? (y - yes, other button - no)
n
Game is over
```

Finished rounds are stored as a JSON array in `results.json`. Each new game is appended to that array:

```json
[
  {
    "date": "2026-10-23T21:35:20.666989098+03:00",
    "outcome": "win",
    "attempts": 4
  }
]
```
