package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const VersionManifestUrl = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type Latest struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

type Version struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"ReleaseTime"`
}

type VersionManifest struct {
	Latest   Latest    `json:"latest"`
	Versions []Version `json:"versions"`
}

func downloadJSON(url string) []byte {
	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Panic(readErr)
	}

	return body
}

func newCli() *cli.App {
	versionManifest := &VersionManifest{}

	app := &cli.App{
		Commands: commands(versionManifest),
		Before: func(c *cli.Context) error {
			body := downloadJSON(VersionManifestUrl)
			jsonErr := json.Unmarshal(body, versionManifest)
			if jsonErr != nil {
				log.Fatal(jsonErr)
			}

			fmt.Println("Download a list", versionManifest.Latest)
			return nil
		},
	}

	return app
}

func commands(manifest *VersionManifest) []*cli.Command {
	return []*cli.Command{
		{
			Name:    "available-versions",
			Aliases: []string{"av"},
			Usage:   "list of available minecraft versions",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "type",
					Usage: "Type of release [release, snapshot, old_alpha, old_beta]",
					Value: "release",
				},
			},
			Action: func(c *cli.Context) error {
				fmt.Println("The latest release:  ", manifest.Latest.Release)
				fmt.Println("The latest snapshot: ", manifest.Latest.Snapshot)
				fmt.Println("Versions:")
				for _, v := range manifest.Versions {
					if t := c.String("type"); len(t) > 0 && t != v.Type {
						continue
					}

					fmt.Printf("[%v] %s - %s\n", v.ReleaseTime, v.Type, v.ID)
				}
				return nil
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "options for servers",
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "add a new server",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Value: "Minecraft server"},
						&cli.StringFlag{Name: "version", Aliases: []string{"v"}, Value: manifest.Latest.Release},
					},
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						fmt.Println("message:", c.String("name"))
						fmt.Println("message:", c.String("version"))
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing server",
					Flags: []cli.Flag{
						&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
					},
					Action: func(c *cli.Context) error {
						fmt.Println("removed server: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}
}
