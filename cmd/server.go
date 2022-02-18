package cmd

import (
	"fmt"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/app"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/http"
	"github.com/OCharnyshevich/Awesome-Minecraft-Server-Wrapper/minecraft"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func Server(config *app.Config, versionManifest *minecraft.VersionManifest, wrappers minecraft.Wrappers) *cli.Command {
	globalPath := fmt.Sprintf("%s/%s", config.RootDir, config.ServerDir)

	return &cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "options for servers",
		Subcommands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a new server",
				Flags: []cli.Flag{
					&cli.StringFlag{Required: true, Name: "name", Aliases: []string{"n"}, Value: "Minecraft server"},
					&cli.StringFlag{Required: true, Name: "version", Aliases: []string{"v"}, Value: "1.18.1"},
					&cli.BoolFlag{Name: "overwrite", Aliases: []string{"o"}, Value: false},
				},
				Action: func(c *cli.Context) error {
					name := c.String("name")
					version := c.String("version")
					overwrite := c.Bool("overwrite")
					versionManifest, err := versionManifest.GetDetails(version)
					if err != nil {
						log.Fatal(err)
					}

					//TODO: create wrapper, attach all function to wrapper
					dir, err := createDir(fmt.Sprintf("%s/%s", globalPath, name))

					if err != nil {
						if !overwrite && os.IsExist(err) {
							log.Fatalf("directory '%s' already exist", dir)
						} else if !os.IsExist(err) {
							log.Fatal(err)
						}
					}

					if err := http.DownloadFile(fmt.Sprintf("%s/%s", dir, config.JarName), versionManifest.Downloads["server"].URL); err != nil {
						log.Fatal(err)
					}

					if err := createELU(dir); err != nil {
						log.Fatal(err)
					}

					fmt.Println("Name:", name)
					fmt.Println("Version:", version)
					fmt.Println("Directory:", dir)
					fmt.Println("URL to download server: ", versionManifest.Downloads["server"].URL)
					fmt.Println("Java version: ", versionManifest.JavaVersion.MajorVersion)

					wpr := minecraft.NewDefaultWrapper(config, name, version, "/usr/lib/jvm/java-17-openjdk/bin/java") //TODO: make chooser to javaVersion
					if err := wrappers.AddWrapper(wpr); err != nil {
						return err
					}
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
			{
				Name:  "run",
				Usage: "run a server",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "name", Aliases: []string{"n"}},
				},
				Action: func(c *cli.Context) error {
					name := c.String("name")

					fmt.Println("Global path: ", globalPath)
					fmt.Println("The server name: ", name)

					wpr, err := wrappers.GetWrapper(name)
					if err != nil {
						return err
					}

					//defer wpr.Stop()

					if err := wpr.Start(); err != nil {
						log.Fatal(err)
						return nil
					}

					//go wrappers.HookStdin(name)

					go func() {
						for {
							select {
							case ev, ok := <-wpr.GameEvents():
								if !ok {
									log.Println("Game events channel closed", ev.String())
									return
								}

								log.Println("events", ev.String())
							}
						}
					}()

					return nil
				},
			},
		},
	}
}

func createDir(path string) (string, error) {
	err := os.Mkdir(path, 0700)

	return path, err
}

func createELU(path string) error {
	out, err := os.Create(fmt.Sprintf("%s/eula.txt", path)) // Create the file
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.WriteString("eula=true")
	return err
}
