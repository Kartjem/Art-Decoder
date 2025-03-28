package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func checkErrors(input string) error {
	re := regexp.MustCompile(`\[([0-9]+)\s([^\]]+)\]`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		if len(match) != 3 {
			return fmt.Errorf("Error")
		}
		if _, err := strconv.Atoi(match[1]); err != nil {
			return fmt.Errorf("Error")
		}
		if match[2] == "" {
			return fmt.Errorf("Error")
		}
	}

	if !isBalanced(input) {
		return fmt.Errorf("Error")
	}

	return nil
}
//
func main() {
	multi := flag.Bool("multiline", false, "enable multiline input")
	encode := flag.Bool("encode", false, "enable encode mode")
	flag.Parse()

	if *multi || *encode {
		var inputLines []string

		if *multi {
			inputLines = readMultiLineInput()
		} else {
			inputLines = flag.Args()
		}

		input := strings.Join(inputLines, "\n")

		if err := checkErrors(input); err != nil {
			fmt.Println(err)
			return
		}

		if *encode {
			encodeInput(input)
		} else {
			processInput(input)
		}
	} else if len(flag.Args()) > 0 {
		input := flag.Args()[0]

		if err := checkErrors(input); err != nil {
			fmt.Println(err)
			return
		}

		if *encode {
			encodeInput(input)
		} else {
			processInput(input)
		}
	} else {
		fmt.Println("Usage: go run main.go (-encode) (-multiline) (-input)")
	}
}

func readMultiLineInput() []string {
	fmt.Println("Enter your text (press 'ENTER' on a new line to finish):")
	scanner := bufio.NewScanner(os.Stdin)
	var inputLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		inputLines = append(inputLines, line)
	}
	if scanner.Err() != nil {
		return []string{}
	}
	return inputLines
}

func processInput(input string) {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if !isBalanced(line) {
			fmt.Println("Error")
			return
		}
		if err := processLine(line); err != nil {
			fmt.Println("Error")
			return
		}
	}
}

func encodeInput(input string) {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		encodedLine := encodeLine(line)
		fmt.Println(encodedLine)
	}
}

func encodeLine(input string) string {
	var result strings.Builder
	count := 1
	prevChar := input[0]

	for i := 1; i < len(input); i++ {
		if input[i] == prevChar {
			count++
		} else {
			if count > 1 {
				result.WriteString(fmt.Sprintf("[%d %s]", count, string(prevChar)))
			} else {
				result.WriteString(string(prevChar))
			}
			prevChar = input[i]
			count = 1
		}
	}

	if count > 1 {
		result.WriteString(fmt.Sprintf("[%d %s]", count, string(prevChar)))
	} else {
		result.WriteString(string(prevChar))
	}

	return result.String()
}

func processLine(line string) error {
	var result strings.Builder
	i := 0

	for i < len(line) {
		if line[i] == '[' {
			countStart := i + 1
			for line[i] != ']' {
				i++
				if i == len(line) {
					return fmt.Errorf("Error")
				}
			}
			countAndChars := line[countStart:i]
			parts := splitCountAndChars(countAndChars)

			if len(parts) != 2 || parts[1] == "" || !startsWithNumber(parts[0]) {
				return fmt.Errorf("Error")
			}

			count, err := strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("Error")
			}
			chars := parts[1]

			for j := 0; j < count; j++ {
				result.WriteString(chars)
			}
			i++
		} else {
			result.WriteByte(line[i])
			i++
		}
	}

	fmt.Println(result.String())
	return nil
}

func isBalanced(s string) bool {
	count := 0
	for _, ch := range s {
		if ch == '[' {
			count++
		} else if ch == ']' {
			count--
		}
		if count < 0 {
			return false
		}
	}
	return count == 0
}

func splitCountAndChars(s string) []string {
	parts := strings.SplitN(s, " ", 2)
	if len(parts) == 1 {
		return parts
	}
	return parts
}

func startsWithNumber(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
