package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/new-timlieberman/3f1_CDC-Study/internal/quiz"
)

func main() {
	modules, err := quiz.LoadModules()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load modules:", err)
		os.Exit(1)
	}
	if len(modules) == 0 {
		fmt.Fprintln(os.Stderr, "no study modules found")
		os.Exit(1)
	}

	statsPath, err := quiz.StatsPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to resolve stats file location:", err)
		os.Exit(1)
	}
	stats, err := quiz.LoadStats(statsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load saved progress:", err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("3F1 CDC Study Helper")
	for {
		questions := selectQuestions(reader, modules, stats)
		if questions == nil {
			fmt.Println("Good luck on your CDC test!")
			return
		}
		quiz.Run(reader, questions, stats, statsPath)

		fmt.Print("\nStudy another set? (y/n): ")
		line, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(line)) != "y" {
			fmt.Println("Good luck on your CDC test!")
			return
		}
	}
}

func selectQuestions(reader *bufio.Reader, modules []quiz.Module, stats quiz.Stats) []quiz.Question {
	var all []quiz.Question
	for _, m := range modules {
		all = append(all, m.Questions...)
	}
	due := dueQuestions(all, stats)

	fmt.Println("\nAvailable modules:")
	for _, m := range modules {
		fmt.Printf("  %d) %s (%d questions)\n", m.ID, m.Name, len(m.Questions))
	}
	fmt.Println("  0) All modules")
	fmt.Printf("  r) Review due questions (%d due)\n", len(due))
	fmt.Println("  q) Quit")

	for {
		fmt.Print("\nChoose a module: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil
		}
		line = strings.TrimSpace(line)
		if strings.EqualFold(line, "q") {
			return nil
		}
		if strings.EqualFold(line, "r") {
			if len(due) == 0 {
				fmt.Println("No questions due for review right now.")
				continue
			}
			return due
		}

		choice, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("Please enter a valid number.")
			continue
		}

		if choice == 0 {
			return all
		}

		for _, m := range modules {
			if m.ID == choice {
				return m.Questions
			}
		}
		fmt.Println("No module with that number, try again.")
	}
}

func dueQuestions(all []quiz.Question, stats quiz.Stats) []quiz.Question {
	var due []quiz.Question
	for _, q := range all {
		if stats.Due(q) {
			due = append(due, q)
		}
	}
	return due
}
