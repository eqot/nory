package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"

	"artifact"
	"gradle"
)

func main() {
	artifact := &artifact.Maven{}

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search artifact",
			Action: func(c *cli.Context) error {
				arts, err := artifact.Find(c.Args().First())
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
				arts, err := artifact.Find(c.Args().First())
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
				arts := gradle.GetArtifacts()

				numOfArts := 0
				var latestArts []string
				for _, art := range arts {
					version := artifact.GetLatestVersion(art)

					if strings.Split(art, ":")[2] < version {
						version = highlight(version)
						numOfArts++
					}

					if version != "" {
						latestArts = append(latestArts, art+":"+version)
					} else {
						latestArts = append(latestArts, art+":")
					}
				}

				renderResult(latestArts)

				if numOfArts > 0 {
					fmt.Println("Artifact(s) can be updated.")
				} else {
					fmt.Println("All artifact(s) are latest.")
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func renderResult(arts []string) {
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

func highlight(text string) string {
	return color.New(color.FgGreen).SprintFunc()(text)
}
