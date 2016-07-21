package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

				fmt.Println("\u2713 Successfully installed.")

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func renderResult(arts []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Group Id", "Artifact Id", "Version"})

	for _, art := range arts {
		table.Append(strings.Split(art, ":"))
	}

	table.Render()
}
