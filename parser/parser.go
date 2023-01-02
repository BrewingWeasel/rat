package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/BurntSushi/toml"
)

type DuplicateRules struct {
	RemoveString     string
	RemoveContinued  []string
	FirstOnly        bool
	RepeatOnTouching bool
	ExactMatch       bool
	KeepEnd          bool
}

type Rules struct {
	Merge            []string
	RemoveDuplicates []DuplicateRules
}

type Location struct {
	file int
	line int
}

func compareBytes(bytes1, bytes2 []byte) bool {
	if len(bytes1) != len(bytes2) {
		return false
	}
	for i, val := range bytes1 {
		if val != bytes2[i] {
			return false
		}
	}
	return true
}

func GenerateLoc(fileNum, lineNum int) Location {
	return Location{file: fileNum, line: lineNum}
}

func GenerateRules(ruleType string) (Rules, error) {
	var rules Rules
	data, err := os.ReadFile(fmt.Sprintf("rules/%s.toml", ruleType))
	if err != nil {
		return rules, err
	}
	toml.Decode(string(data), &rules)
	return rules, nil
}

func UseRules(rules Rules, input [][]byte) (map[Location]string, error) {
	changes := map[Location]string{}
	for _, remove := range rules.RemoveDuplicates { // TODO Change everything so it has to loop less
		lastMatchVal := []byte{}
		lastMatch := false
		previousMatches := []Location{}
		for fileIndex, file := range input {
			bytesReader := bytes.NewReader(file)
			bufReader := bufio.NewReader(bytesReader)
			lineNum := 0
			for {
				search, err := bufReader.ReadBytes('\n')
				if err != nil {
					break
				}
				match, err := regexp.Match(remove.RemoveString, search)
				if match && err == nil {
					newMatch := true
					for _, newRemoveString := range remove.RemoveContinued {
						nextLine, err := bufReader.Peek('\n') // TODO Make it work for more than one...
						if err != nil {
							break
						}
						newMatch, err = regexp.Match(newRemoveString, nextLine)
						if err != nil {
							return changes, err
						}
						if !newMatch {
							break
						}
					}
					// TODO: Clean up
					if newMatch {
						loc := Location{file: fileIndex, line: lineNum}
						if lastMatch && !remove.ExactMatch {
							if remove.KeepEnd {
								oldLoc := previousMatches[len(previousMatches)-1]
								changes[oldLoc] = ""
								for i := 0; i <= len(remove.RemoveContinued); i++ {
									fmt.Print(i)
									changes[Location{line: oldLoc.line + i, file: oldLoc.file}] = ""
								}
							} else {
								changes[loc] = ""
								for i := 0; i <= len(remove.RemoveContinued); i++ {
									fmt.Print(i)
									changes[Location{line: loc.line + i, file: loc.file}] = ""
								}
							}
						} else if compareBytes(lastMatchVal, search) {
							changes[loc] = ""
						}
						lastMatch = true
						lastMatchVal = search
						previousMatches = append(previousMatches, loc)
					}
				} else if err != nil {
					return changes, err
				}
				if remove.FirstOnly || !match && remove.RepeatOnTouching {
					break
				}
				lineNum++
			}
		}
	}
	return changes, nil
}
