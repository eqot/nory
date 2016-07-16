package main

import (
	"log"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"

	"artifact"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search artifact",
			Action: func(c *cli.Context) error {
				artifact := &artifact.Maven{}
				arts, err := artifact.Find(c.Args().First())
				if err != nil {
					log.Fatal(err)
				}

				renderResult(arts)

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
