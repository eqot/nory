package main

import (
	"fmt"
	"log"
	"os"
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
				arts := gradle.GetArtifacts()

				var wg sync.WaitGroup
				ch := make(chan string, 64)

				for _, art := range arts {
					wg.Add(1)

					go func(art2 string) {
						defer wg.Done()

						artifactRepo := &artifact.Maven{}

						latestArt, _ := artifactRepo.GetLatestVersion(art2)
						if latestArt == "" {
							return
						}

						latestVersion := artifact.GetVersion(latestArt)

						if artifact.GetVersion(art2) < latestVersion {
							latestVersion = highlight(latestVersion)
						}

						ch <- art2 + ":" + latestVersion
					}(art)
				}

				wg.Wait()

				var latestArts []string
				for len(ch) > 0 {
					latestArt := <-ch
					latestArts = append(latestArts, latestArt)
				}

				renderResult(latestArts)

				if len(latestArts) > 0 {
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
				arts := gradle.GetArtifacts()

				var latestArts []string
				for _, art := range arts {
					latestArt, _ := artifactRepo.GetLatestVersion(art)
					if latestArt == "" {
						continue
					}

					latestVersion := artifact.GetVersion(latestArt)

					if artifact.GetVersion(art) < latestVersion {
						gradle.Add(latestArt)

						latestVersion = highlight(latestVersion)
						latestArts = append(latestArts, art+":"+latestVersion)
					}
				}

				renderResult(latestArts)

				if len(latestArts) > 0 {
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

func highlight(text string) string {
	return color.New(color.FgGreen).SprintFunc()(text)
}
