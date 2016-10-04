package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

// ProspectorMultiLineConfig ...
type ProspectorMultiLineConfig struct {
	Pattern string `yaml:"pattern"`
	Negate  bool   `yaml:"negate"`
	Match   string `yaml:"match"`
}

// ProspectorConfig ...
type ProspectorConfig struct {
	Paths           []string                  `yaml:"paths"`
	FieldsUnderRoot bool                      `yaml:"fields_under_root"`
	IgnoreOlder     string                    `yaml:"ignore_older"`
	Fields          map[string]string         `yaml:"fields"`
	MultiLine       ProspectorMultiLineConfig `yaml:"multiline"`
}

// Prospectors ...
type Prospectors struct {
	Prospectors []ProspectorConfig `yaml:"prospectors"`
}

// FilebeatConfig ...
type FilebeatConfig struct {
	Filebeat Prospectors `yaml:"filebeat"`
}

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
	yamlConfig := flag.String("y", "", "Filebeat prospector yaml config file (overrides the pattern/negate options)")

	flag.Parse()

	fmt.Println("")

	if *v {
		fmt.Printf("Version %s (commit: %s, %s)\n", version, commitHash, buildDate)
		os.Exit(0)
	}

	// Choose yaml config first if it was specified
	if *yamlConfig != "" {
		conf, err := loadYamlConfig(*yamlConfig)
		if err != nil {
			exitWithMessage("ERROR", fmt.Sprintf("Problem with yaml config: %s", err), true)
		}
		*pattern = conf.Filebeat.Prospectors[0].MultiLine.Pattern
		*negate = conf.Filebeat.Prospectors[0].MultiLine.Negate
		*file = conf.Filebeat.Prospectors[0].Paths[0]
	}

	if *file == "" {
		exitWithMessage("ERROR", "Must specify a file name.", true)
	}

	if *pattern == "" {
		exitWithMessage("ERROR", "Must specify a pattern.", true)
	}

	content, err := ioutil.ReadFile(*file)
	if err != nil {
		exitWithMessage("ERROR", fmt.Sprintf("Could not read file: %s", err), true)
	}

	regex, err := regexp.Compile(*pattern)
	if err != nil {
		exitWithMessage("ERROR", fmt.Sprintf("Failed to compile pattern: %s", err), false)
	}

	lines := strings.Split(string(content), "\n")

	if string(content) == "" {
		exitWithMessage("WARNING", "Sample string contents is empty.", false)
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

// Parse is a function that unmarshals the specified yaml config file
func (c *FilebeatConfig) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	if len(c.Filebeat.Prospectors) == 0 {
		return errors.New("Must have at least one prospector config!")
	}

	return nil
}

func loadYamlConfig(filname string) (*FilebeatConfig, error) {

	content, err := ioutil.ReadFile(filname)
	if err != nil {
		return nil, err
	}

	var config FilebeatConfig
	if err := config.Parse(content); err != nil {
		return nil, err
	}

	return &config, nil
}

func exitWithMessage(level string, msg string, showUsage bool) {
	fmt.Printf("[%s] %s\n", level, msg)
	if showUsage {
		flag.PrintDefaults()
	}
	os.Exit(1)
}
