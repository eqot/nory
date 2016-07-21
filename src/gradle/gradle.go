package gradle

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"artifact"
)

const filename string = "./app/build.gradle"
const tmpFilename string = "./app/build.gradle.tmp"

const (
	stateNormal = iota
	stateStartOfDependencies
	stateDependency
	stateEndOfDependencies
)

func Add(art string) {
	var buffer string
	var isAlreadyAdded bool
	Parse(func(state int, line string) {
		switch state {
		case stateDependency:
			art2 := getArt(line)
			if art2 != "" && artifact.IsSameArtifact(art, art2) {
				buffer = injectDependency(buffer, artifact.GetLatest(art, art2))

				isAlreadyAdded = true

				return
			}

		case stateEndOfDependencies:
			if !isAlreadyAdded {
				buffer = injectDependency(buffer, art)
			}
		}

		buffer = fmt.Sprintf("%s%s\n", buffer, line)
	})

	fo, err := os.Create(tmpFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer fo.Close()

	so := bufio.NewWriter(fo)
	defer so.Flush()

	so.WriteString(buffer)

	if err := os.Rename(tmpFilename, filename); err != nil {
		log.Fatal(err)
	}
}

func injectDependency(buffer, art string) string {
	return fmt.Sprintf("%s    compile '%s'\n", buffer, art)
}

func getArt(line string) string {
	re := regexp.MustCompile(`\s*compile\s+'(.+:.+:.+)'`)
	match := re.FindStringSubmatch(line)
	if len(match) > 0 {
		return match[1]
	}

	return ""
}

func Parse(callback func(int, string)) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var state int
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		if state == stateStartOfDependencies {
			state = stateDependency
		}

		if line == "dependencies {" {
			state = stateStartOfDependencies
		}

		if state == stateDependency && line == "}" {
			state = stateEndOfDependencies
		}

		callback(state, line)
	}
}

func Read() {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	re := regexp.MustCompile(`(\s*)compile\s'(.+):(.+):(.+)'`)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		match := re.FindStringSubmatch(line)

		if len(match) == 0 {
			fmt.Println(line)
		} else {
			fmt.Print(match[1])
			fmt.Println(match[2:])
		}
	}
}

func Parse2() ([]string, error) {
	var result []string
	var lines string

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var isDependencies bool

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()

		if line == "dependencies {" {
			isDependencies = true

			result = append(result, lines)
			lines = ""
		}

		if isDependencies && line == "}" {
			isDependencies = false

			result = append(result, lines)
			lines = ""
		}

		lines = fmt.Sprintf("%s%s\n", lines, line)
	}

	result = append(result, lines)

	return result, nil
}
