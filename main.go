package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Stack ...
type Stack struct {
	top  *Element
	size int
}

// Element ...
type Element struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next  *Element
}

// Len function returns the stack's length
func (s *Stack) Len() int {
	return s.size
}

// Push function a new element onto the stack
func (s *Stack) Push(value interface{}) {
	s.top = &Element{value, s.top}
	s.size++
}

// Pop function removes the top element from the stack and return it's value or nil if the stack is empty
func (s *Stack) Pop() (value interface{}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

// Peek function returns the top element value from the stack without modifying the stack
func (s *Stack) Peek() (value interface{}) {
	return s.top.value
}

// Reset function resets the stack by setting the size to 0 and the top element to nil
func (s *Stack) Reset() {
	s.top = nil
	s.size = 0
}

var (
	buildDate  string
	version    string
	commitHash string
)

func main() {

	v := flag.Bool("v", false, "prints current version and exits")
	pattern := flag.String("p", "", "Multi-line regex pattern")
	negate := flag.Bool("n", true, "Negate the pattern matching")
	file := flag.String("f", "", "File containing multi-line string")
	flag.Parse()

	fmt.Println("")

	if *v {
		fmt.Printf("Version %s (commit: %s, %s)\n", version, commitHash, buildDate)
		os.Exit(0)
	}

	if *file == "" {
		fmt.Println("[ERROR] Must specify a file name!\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println("[ERROR] Could not read file: ", err, "\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *pattern == "" {
		fmt.Println("[ERROR] Must specify a pattern!\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	regex, err := regexp.Compile(*pattern)
	if err != nil {
		fmt.Println("Failed to compile pattern: ", err, "\n")
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")

	if string(content) == "" {
		fmt.Println("[WARNING] Sample string contents is empty!")
		os.Exit(1)
	}

	fmt.Println("Pattern Match?\t\tString\n--------------------------------------------")

	totalFullMatches := 0

	stackHead := new(Stack)
	stackTail := new(Stack)

	for _, line := range lines {

		matches := regex.MatchString(line)
		if *negate {
			matches = !matches
		}

		if stackHead.Len() == 0 {

			stackHead.Push(matches)

		} else {

			// Check if the item to be pushed is the same as what's in the current head stack
			// If so, then restet both stacks
			if matches == stackHead.Peek().(bool) {

				// If the tail stack has one or more items, count this as a match before reseting both stacks
				if stackTail.Len() >= 1 {
					totalFullMatches += 1
				}
				stackHead.Reset()
				stackHead.Push(matches)
				stackTail.Reset()

			} else {
				// Otherwise, push the match result onto the tail stack
				stackTail.Push(matches)
			}
		}

		fmt.Printf("%v\t\t%v\n", matches, line)

	}

	fmt.Println("")
	fmt.Println("-------------------------")
	fmt.Printf("Total Matches: %d\n", totalFullMatches)
	fmt.Println("-------------------------")

}
