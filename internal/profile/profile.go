// Package profile loads everything in data/ that drives the card and README.
package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Links struct {
	GitHub   Link `json:"github"`
	LinkedIn Link `json:"linkedin"`
}

// Item is one "├ Key → Value" line under a section.
type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Section is one "KEY → Value" block on the card. Set Blank for a spacer line
// and Stats for the live GitHub stats line.
type Section struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Items []Item `json:"items"`
	Blank bool   `json:"blank"`
	Stats bool   `json:"stats"`
}

// Profile is the identity/presentation data driving the card and README.
type Profile struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName"`
	Host        string    `json:"host"`
	Title       string    `json:"title"`
	Email       string    `json:"email"`
	Bio         []string  `json:"bio"`
	Sections    []Section `json:"sections"`
	Links       Links     `json:"links"`
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

func GetProfile() (Profile, error) {
	var p Profile
	return p, readJSON("data/profile.json", &p)
}

// GetPortrait reads the ASCII portrait, dropping leading/trailing blank lines.
func GetPortrait() ([]string, error) {
	data, err := os.ReadFile("data/portrait.txt")
	if err != nil {
		return nil, fmt.Errorf("read data/portrait.txt: %w", err)
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("data/portrait.txt is empty")
	}
	return lines, nil
}
