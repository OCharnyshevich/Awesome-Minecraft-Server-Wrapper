package cmd

import (
	"fmt"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/minecraft"
	"github.com/urfave/cli/v2"
)

func AvailableVersions(manifest *minecraft.VersionManifest) *cli.Command {
	return &cli.Command{
		Name:    "versions",
		Aliases: []string{"v"},
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
	}
}
