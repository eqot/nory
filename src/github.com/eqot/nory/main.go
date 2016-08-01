package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"

	"github.com/eqot/nory/lib/artifact"
	"github.com/eqot/nory/lib/gradle"
)

func main() {
	artifactRepo := &artifact.Maven{}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search artifact",
			Action: func(c *cli.Context) error {
				arts, err := artifactRepo.Find(c.Args().First())
				if err != nil {
					log.Fatal(err)
				}

				renderResult(arts)

				return nil
			},
		},

		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install artifact",
			Action: func(c *cli.Context) error {
				arts, err := artifactRepo.Find(c.Args().First())
				if err != nil {
					log.Fatal(err)
				}

				gradle.Add(arts[0])

				renderResult([]string{arts[0]})

				fmt.Printf("%s Successfully installed.\n", highlight("\u2713"))

				return nil
			},
		},

		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "check if artifacts can be updated",
			Action: func(c *cli.Context) error {
				updatableArts := getUpdatableArtifacts()
				updatableArts = highlightUpdatableArtifacts(updatableArts)
				renderResult(updatableArts)

				if len(updatableArts) > 0 {
					fmt.Println("Artifact(s) can be updated.")
				} else {
					fmt.Println("No outdated artifacts.")
				}

				return nil
			},
		},

		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update outdated artifacts if exists",
			Action: func(c *cli.Context) error {
				updatableArts := getUpdatableArtifacts()

				for _, art := range updatableArts {
					arts := artifact.Split(art)
					latestArt := arts[0] + ":" + arts[1] + ":" + arts[3]

					gradle.Add(latestArt)
				}

				updatableArts = highlightUpdatableArtifacts(updatableArts)
				renderResult(updatableArts)

				if len(updatableArts) > 0 {
					fmt.Printf("%s Successfully updated.\n", highlight("\u2713"))
				} else {
					fmt.Println("No outdated artifacts.")
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func getUpdatableArtifacts() []string {
	arts := gradle.GetArtifacts()

	ch := make(chan string, 64)
	var wg sync.WaitGroup

	for _, art := range arts {
		wg.Add(1)

		go artifact.GetArtifactWithLatestVersion(art, ch, &wg)
	}

	wg.Wait()

	var updatableArts []string
	for len(ch) > 0 {
		art := <-ch

		if artifact.IsUpdatable(art) {
			updatableArts = append(updatableArts, art)
		}
	}

	return updatableArts
}

func highlightUpdatableArtifacts(arts []string) []string {
	var result []string

	for _, art := range arts {
		if artifact.IsUpdatable(art) {
			art = highlightColumn(art, 3)
		}

		result = append(result, art)
	}

	return result
}

func highlightColumn(art string, index int) string {
	arts := artifact.Split(art)

	arts[index] = highlight(arts[index])

	return strings.Join(arts, ":")
}

func highlight(text string) string {
	return color.New(color.FgGreen).SprintFunc()(text)
}

func renderResult(arts []string) {
	if arts == nil || len(arts) == 0 {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)

	title := color.New(color.FgCyan).SprintFunc()

	if len(artifact.Split(arts[0])) == 4 {
		table.SetHeader([]string{
			title("Group Id"), title("Artifact Id"), title("Current"), title("Latest"),
		})
	} else {
		table.SetHeader([]string{
			title("Group Id"), title("Artifact Id"), title("Version"),
		})
	}

	table.SetAutoFormatHeaders(false)

	for _, art := range arts {
		table.Append(artifact.Split(art))
	}

	table.Render()
}
