package main

import (
	"errors"
	"fmt"
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

func line(size int) {
	fmt.Println(strings.Repeat("-", size))
}

func rightjust(s string, n int, fill string) string {
	if len(s) > n {
		return s
	}
	return strings.Repeat(fill, (n-len(s))) + s
}

func leftjust(s string, n int, fill string) string {
	if len(s) > n {
		return s
	}
	return s + strings.Repeat(fill, (n-len(s)))
}

func center(s string, n int, fill string) string {
	div := (n - len(s)) / 2
	return strings.Repeat(fill, div) + s + strings.Repeat(fill, div)
}
