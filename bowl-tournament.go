package main

// @author: Ashley Kitson
// @date:   20/05/2024
// @project: bowl-tournament
// @file:   bowl-tournament.go
// @license: MIT

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"github.com/thatisuday/commando"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	PlantumlPath = "/usr/share/bowl-tournament/plantuml.jar"
)

type Match struct {
	Round           int
	Match           string
	PlayBy          *time.Time
	PlayOnDate      *time.Time
	PlayOnTime      *string
	HomeParticipant string
	AwayParticipant string
	NextMatch       *string
}
type Matches map[string]*Match

type input []string

//go:embed VERSION
var version string

func init() {
	commando.
		SetExecutableName("bowl-tournament").
		SetDescription("Bowl tournament utility").
		SetVersion(version)
	commando.
		Register(nil).
		SetDescription("create PNG file of csv input for bowls tounament").
		AddArgument("input", "input csv file", "").
		AddArgument("tournament", "name of tournament", "").
		AddArgument("game", "name of game", "").
		AddFlag("groupdates,g", "group dates in the round data", commando.Bool, false).
		AddFlag("theme,t", "the Plantuml theme to use. See https://plantuml.com/theme", commando.String, "_none_").
		AddFlag("output,o", "output png file", commando.Bool, false).
		AddFlag("glob,l", "find and process all csv files. 'input' argument is the directory path to process. 'game' argument is regex pattern to extract game name from file name", commando.Bool, false).
		AddFlag("plantumlpath,p", "path to plantuml.jar if not using supplied jar", commando.String, PlantumlPath).
		SetAction(command)

}

func command(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	glob, err := flags["glob"].GetBool()
	if err != nil {
		panic(err)
	}
	inputs := make(input, 0)
	walkDir := func(root string) error {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && path != root {
				return filepath.SkipDir
			}
			if strings.HasSuffix(path, ".csv") {
				inputs = append(inputs, path)
				return nil
			}
			return nil
		})
		return nil
	}
	if glob {
		err := walkDir(args["input"].Value)
		if err != nil {
			panic(err)
		}
	} else {
		inputs = append(inputs, args["input"].Value)
		if _, err := os.Stat(inputs[0]); err != nil {
			panic(err)
		}
	}

	for _, inp := range inputs {
		file, err := os.Open(inp)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			panic(err)
		}
		//remove header
		records = records[1:]

		matches := make(Matches)
		numRounds := 0
		for _, record := range records {
			r, err := strconv.Atoi(record[0])
			if err != nil {
				panic(err)
			}
			if r > numRounds {
				numRounds = r
			}
			var t1, t2 *time.Time
			if record[2] != "" {
				tx, err := time.Parse("02/01/06", record[2])
				if err != nil {
					panic(err)
				}
				t1 = &tx
			}
			if record[3] != "" {
				tx, err := time.Parse("02/01/06", record[3])
				if err != nil {
					panic(err)
				}
				t2 = &tx
			}
			var t3 *string
			if record[4] != "" {
				t3 = &record[4]
			}
			match := Match{
				Round:           r,
				Match:           record[1],
				PlayBy:          t1,
				PlayOnDate:      t2,
				PlayOnTime:      t3,
				HomeParticipant: record[5],
				AwayParticipant: record[6],
			}
			matches[match.Match] = &match
		}

		//work out the match hierarchy
		for _, match := range matches {
			if strings.Contains(match.HomeParticipant, "Winner of") {
				prevMatch := strings.Split(match.HomeParticipant, " ")[2]
				matches[prevMatch].NextMatch = &match.Match
			}
			if strings.Contains(match.AwayParticipant, "Winner of") {
				prevMatch := strings.Split(match.AwayParticipant, " ")[2]
				matches[prevMatch].NextMatch = &match.Match
			}
		}

		groupDates, err := flags["groupdates"].GetBool()
		if err != nil {
			panic(err)
		}

		theme, _ := flags["theme"].GetString()
		var game string
		if glob {
			re := regexp.MustCompile(args["game"].Value)
			g := re.FindStringSubmatch(strings.TrimSuffix(filepath.Base(inp), filepath.Ext(inp)))
			game = g[1]
		} else {
			game = args["game"].Value
		}
		puml := createPuml(matches, numRounds, groupDates, args["tournament"].Value, game, theme)
		outfile := strings.TrimSuffix(inp, filepath.Ext(inp)) + ".puml"
		f, err := os.Create(outfile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.WriteString(puml)
		if err != nil {
			panic(err)
		}

		createPng, err := flags["output"].GetBool()
		if err != nil {
			panic(err)
		}
		if createPng {
			var pPath string
			if args["plantumlpath"].Value == "" {
				pPath = PlantumlPath
			} else {
				pPath = args["plantumlpath"].Value
			}
			createPngFile(inp, pPath)
		}
	}
}

func createPuml(matches Matches, numRounds int, groupDates bool, tournament, game, theme string) string {

	//find the optimal minimal width
	const lex = 7.57757
	numChars := 0
	for _, m := range matches {
		if len(m.HomeParticipant) > numChars {
			numChars = len(m.HomeParticipant)
		}
		if len(m.AwayParticipant) > numChars {
			numChars = len(m.AwayParticipant)
		}
	}
	minWidth := int(lex * float64(numChars))

	puml := fmt.Sprintf(`@startuml
title %s\n%s
hide empty methods
hide circle
left to right direction
skinparam minClassWidth %d
!theme %s
`, tournament, game, minWidth, theme)

	// find out how many matches in each level
	numMatchesPerLevel := make(map[int]int)
	for r := 1; r <= numRounds; r++ {
		for _, m := range matches {
			if m.Round == r {
				numMatchesPerLevel[r]++
			}
		}
	}

	// for each round create a package
	for r := 1; r <= numRounds; r++ {
		//work out round name
		var pkgName string
		if numMatchesPerLevel[r] == 1 && r == numRounds {
			pkgName = "Finals"
		} else if numMatchesPerLevel[r] == 2 && r == numRounds-1 {
			pkgName = "Semi Finals"
		} else {
			pkgName = "Round " + strconv.Itoa(r)
		}

		puml += fmt.Sprintf("package \"%s\" {\n", pkgName)

		var playDateWritten bool
		//create the rounds for the current package
		for _, m := range matches {
			if m.Round == r {
				if m.PlayBy != nil && groupDates && !playDateWritten {
					puml += fmt.Sprintf("\tnote \"Play By: <b>%s</b>\" as n%d\n", m.PlayBy.Format("02/01/06"), r)
					playDateWritten = true
				}
				if m.PlayOnDate != nil && groupDates && !playDateWritten {
					puml += fmt.Sprintf("\tnote \"Play On: <b>%s %s</b>\" as n%d\n", m.PlayOnDate.Format("02/01/06"), *m.PlayOnTime, r)
					playDateWritten = true
				}

				puml += fmt.Sprintf("\tclass %s {\n", m.Match)
				if m.PlayBy != nil && !groupDates {
					puml += fmt.Sprintf("\t\tPlay By: <b>%s</b>\n", m.PlayBy.Format("02/01/06"))
				}
				if m.PlayOnDate != nil && !groupDates {
					puml += fmt.Sprintf("\t\tPlay On: <b>%s %s</b>\n", m.PlayOnDate.Format("02/01/06"), *m.PlayOnTime)
				}
				puml += fmt.Sprintf("\t\t%s\n", m.HomeParticipant)
				if strings.Contains(m.HomeParticipant, "Winner of") {
					puml += "-\n-\n"
				}
				puml += fmt.Sprintf("\t\t%s\n", m.AwayParticipant)
				if strings.Contains(m.AwayParticipant, "Winner of") {
					puml += "-\n-\n"
				}
				puml += "\t}\n"
			}
		}
		playDateWritten = false

		puml += "}\n"
	}
	// add links
	for _, m := range matches {
		if m.NextMatch != nil {
			puml += fmt.Sprintf("%s --> %s\n", m.Match, *m.NextMatch)
		}
	}

	puml += `@enduml`
	return puml
}

func createPngFile(input string, plantumlpath string) {
	absPath, err := filepath.Abs(input)
	if err != nil {
		panic(err)
	}
	outfile := strings.TrimSuffix(absPath, filepath.Ext(absPath)) + ".puml"
	cmd := exec.Command("java", "-jar", plantumlpath, outfile)
	println(cmd.String())
	_, err = cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func main() {
	commando.Parse(nil)
	os.Exit(0)
}
