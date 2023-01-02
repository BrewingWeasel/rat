package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BrewingWeasel/rat/parser"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/urfave/cli/v2"
)

func getTypeFromName(name string) string {
	ending := strings.Split(name, ".")
	if len(ending) == 1 {
		return "text"
	}
	return ending[1]
}

func contains(source map[parser.Location]string, loc parser.Location) bool {
	for i := range source {
		if i == loc {
			return true
		}
	}
	return false
}

func printContents(transformations map[parser.Location]string, fileNum int, contents string, visibleLineNum *int, lines bool, fileType string) error {
	lexer := lexers.Get(fileType)
	if fileType == "text" {
		lexer = lexers.Get("bash")
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	style := styles.Get("fruity")
	formatter := formatters.TTY256

	lineStrings := strings.Split(contents, "\n")

	for lineNum, curLine := range lineStrings {
		loc := parser.GenerateLoc(fileNum, lineNum)
		isChanged := contains(transformations, loc)
		if isChanged {
			// Do check for other changes later lol
		} else {
			*visibleLineNum++
			if lines {
				iterator, _ := lexer.Tokenise(nil, curLine)
				buf := new(bytes.Buffer)
				formatter.Format(buf, style, iterator)
				fmt.Println(*visibleLineNum, buf.String())
			} else {
				fmt.Println(curLine)
			}
		}
	}
	return nil
}

func main() {
	lines := false

	app := &cli.App{
		Name:  "Rat",
		Usage: "A smart version of cat",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "numbers",
				Aliases:     []string{"n"},
				Usage:       "Disable showing line numbers",
				Destination: &lines,
			},
		},
		Action: func(ctx *cli.Context) error {
			lineNum := 0
			files := [][]byte{}

			lastType := ""
			sameType := true
			for i := 0; i < ctx.NArg(); i++ {
				// Check if files are of the same type (can they be smartly concatenated)
				if sameType {
					curType := getTypeFromName(ctx.Args().Get(1))
					if lastType != "" && lastType != curType {
						sameType = false
					} else {
						lastType = curType
					}
				}

				file, err := os.ReadFile(ctx.Args().Get(i))
				if err != nil {
					return err
				}
				files = append(files, file)
			}

			transformations := map[parser.Location]string{}

			if sameType && ctx.NArg() != 1 {
				rules, err := parser.GenerateRules(lastType)
				if err != nil {
					return err
				}
				transformations, err = parser.UseRules(rules, files)
				if err != nil {
					return err
				}
			}

			for i := 0; i < ctx.NArg(); i++ {
				err := printContents(transformations, i, string(files[i]), &lineNum, lines, lastType)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
