package quiz

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var labels = []string{"A", "B", "C", "D", "E", "F"}

type attempt struct {
	question Question
	choices  []string
	answer   int
}

// Run drives an interactive quiz session over questions, looping through
// missed-question retry rounds until the user answers everything correctly
// or declines to retry. Each attempt updates stats (spaced-repetition
// progress), which is persisted to statsPath as it changes.
func Run(reader *bufio.Reader, questions []Question, stats Stats, statsPath string) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	pool := make([]Question, len(questions))
	copy(pool, questions)

	round := 1
	for {
		if round > 1 {
			fmt.Printf("\n--- Retry round %d: %d missed question(s) ---\n", round, len(pool))
		}

		missed := runRound(reader, rng, pool, stats, statsPath)
		if len(missed) == 0 {
			fmt.Println("\nAll questions answered correctly!")
			return
		}

		if !askYesNo(reader, fmt.Sprintf("\nRetry the %d missed question(s)? (y/n): ", len(missed))) {
			return
		}
		pool = missed
		round++
	}
}

func runRound(reader *bufio.Reader, rng *rand.Rand, questions []Question, stats Stats, statsPath string) []Question {
	order := rng.Perm(len(questions))
	correct := 0
	var missed []Question

	for i, idx := range order {
		q := questions[idx]
		att := shuffleChoices(rng, q)

		fmt.Printf("\nQ%d. %s\n", i+1, att.question.Question)
		for j, c := range att.choices {
			fmt.Printf("  %s) %s\n", labels[j], c)
		}

		given := promptAnswer(reader, len(att.choices))
		isCorrect := given == att.answer
		if isCorrect {
			fmt.Println("Correct!")
			correct++
		} else {
			fmt.Printf("Incorrect. Correct answer: %s) %s\n", labels[att.answer], att.choices[att.answer])
			if q.Explanation != "" {
				fmt.Printf("Why: %s\n", q.Explanation)
			}
			missed = append(missed, q)
		}

		stats.Record(q, isCorrect)
		if err := SaveStats(statsPath, stats); err != nil {
			fmt.Fprintln(os.Stderr, "warning: failed to save progress:", err)
		}
	}

	total := len(questions)
	pct := 0.0
	if total > 0 {
		pct = float64(correct) / float64(total) * 100
	}
	fmt.Printf("\nScore: %d/%d (%.0f%%)\n", correct, total, pct)

	return missed
}

func shuffleChoices(rng *rand.Rand, q Question) attempt {
	order := rng.Perm(len(q.Choices))
	choices := make([]string, len(q.Choices))
	newAnswer := 0
	for newIdx, origIdx := range order {
		choices[newIdx] = q.Choices[origIdx]
		if origIdx == q.Answer {
			newAnswer = newIdx
		}
	}
	return attempt{question: q, choices: choices, answer: newAnswer}
}

func promptAnswer(reader *bufio.Reader, numChoices int) int {
	for {
		fmt.Print("Your answer: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return -1
		}
		line = strings.ToUpper(strings.TrimSpace(line))
		if line == "" {
			continue
		}

		for i := 0; i < numChoices && i < len(labels); i++ {
			if line == labels[i] {
				return i
			}
		}
		if n, err := strconv.Atoi(line); err == nil && n >= 1 && n <= numChoices {
			return n - 1
		}

		fmt.Printf("Please answer with a letter (A-%s) or number (1-%d).\n", labels[numChoices-1], numChoices)
	}
}

func askYesNo(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	line = strings.ToLower(strings.TrimSpace(line))
	return line == "y" || line == "yes"
}