package main

import (
	//"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/elastic/beats/libbeat/common/match"
	//"gopkg.in/yaml.v2"
)

var (
	buildDate       string
	version         string
	commitHash      string
	totalMinMatches int
	stackHead       *Stack
	stackTail       *Stack
	matcherTypes    map[int]string
	showMatches     *bool
)

func main() {

	matcherTypes = map[int]string{
		1: "regexp.CompilePosix",
		2: "regexp.Compile",
		3: "match.Matcher",
	}

	v := flag.Bool("v", false, "prints current version and exits")
	matcherType := flag.Int("t", 0, "Pattern backend type: 1=regexp.CompilePosix, 2=regexp.Compile, 3=match.Matcher")
	pattern := flag.String("p", "", "Multi-line regex pattern")
	negate := flag.Bool("n", true, "Negate the pattern matching")
	file := flag.String("f", "", "File containing multi-line string")
	showMatches = flag.Bool("s", false, "Print individual lines and their matching status")
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

	if *matcherType < 1 || *matcherType > 3 {
		exitWithMessage("ERROR", "Pattern type must be a value from 1 to 3", true)
	}

	fileContent, err := ioutil.ReadFile(*file)
	if err != nil {
		exitWithMessage("ERROR", fmt.Sprintf("Could not read file: %s", err), true)
	}
	content := strings.TrimSpace(string(fileContent))

	if content == "" {
		exitWithMessage("WARNING", "Sample string contents is empty.", false)
	}

	stackHead = new(Stack)
	stackTail = new(Stack)

	line(70)
	fmt.Println(center("Filebeat multiline tester", 70, " "))
	line(70)
	fmt.Printf("Using matching backend: %s\n", matcherTypes[*matcherType])
	line(70)

	switch *matcherType {
	case 1, 2:
		MatchOnRegexpCompile(*pattern, *negate, *matcherType, content)
	case 3:
		MatchOnMatchMatcher(*pattern, *negate, content)
	}

	if *showMatches {
		line(70)
	}
	fmt.Printf("Matches Found: %d\n", totalMinMatches)
	line(70)

}

func MatchOnRegexpCompile(pattern string, negate bool, matcherType int, content string) {

	var regex *regexp.Regexp
	var err error
	if matcherType == 1 {
		regex, err = regexp.CompilePOSIX(pattern)
		if err != nil {
			exitWithMessage("ERROR", fmt.Sprintf("Failed to compile pattern: %s", err), false)
		}
	} else {
		regex, err = regexp.Compile(pattern)
		if err != nil {
			exitWithMessage("ERROR", fmt.Sprintf("Failed to compile pattern: %s", err), false)
		}
	}

	lines := strings.Split(content, "\n")

	if *showMatches {
		fmt.Printf("%s| %s\n", center("Match?", 10, " "), leftjust("Text", 40, " "))
		fmt.Println(strings.Repeat("-", 60))
	}

	for _, line := range lines {

		isMatch := regex.MatchString(line)
		if negate {
			isMatch = !isMatch
		}

		if stackHead.Len() == 0 {
			stackHead.Push(isMatch)
			totalMinMatches++
		} else {
			// Check if the item to be pushed is the same as what's in the current head stack
			// If so, then restet both stacks
			if isMatch == stackHead.Peek().(bool) {
				totalMinMatches++
				stackHead.Reset()
				stackHead.Push(isMatch)
				stackTail.Reset()
			} else {
				// Otherwise, push the match result onto the tail stack
				stackTail.Push(isMatch)
			}
		}
		if *showMatches {
			fmt.Printf("%v| %s\n", center(fmt.Sprintf("%v", isMatch), 10, " "), leftjust(line, 40, " "))
		}
	}

}

func MatchOnMatchMatcher(pattern string, negate bool, content string) {

	matcher := match.MustCompile(pattern)
	lines := strings.Split(content, "\n")

	if *showMatches {
		fmt.Printf("%s| %s\n", center("Match?", 10, " "), leftjust("Text", 40, " "))
		fmt.Println(strings.Repeat("-", 60))
	}

	var isMatchString string
	for _, line := range lines {
		isMatch := matcher.MatchString(line)
		if negate {
			isMatch = !isMatch
		}
		if stackHead.Len() == 0 {
			stackHead.Push(isMatch)
			totalMinMatches++
		} else {
			// Check if the item to be pushed is the same as what's in the current head stack
			// If so, then restet both stacks
			if isMatch == stackHead.Peek().(bool) {
				totalMinMatches++
				stackHead.Reset()
				stackHead.Push(isMatch)
				stackTail.Reset()
			} else {
				// Otherwise, push the match result onto the tail stack
				stackTail.Push(isMatch)
			}
		}

		if isMatch {
			isMatchString = "true"
		} else {
			isMatchString = "false"
		}

		if *showMatches {
			fmt.Printf("%v| %s\n", center(isMatchString, 10, " "), leftjust(line, 40, " "))
		}
	}

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
