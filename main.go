package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var version = "dev"

const (
	nestedEnvKey = "__OSSIE_SPAWNED"
	nestedEnvVal = "righto"
)

type contextKey string

func debugNYI(c *cli.Context) {
	fmt.Printf("Args: %#v\n", c.Args())
	fmt.Printf("Flagnames: %#v\n", c.FlagNames())
	fmt.Printf("LocalFlags: %#v\n", c.LocalFlagNames())
	fmt.Printf("Flags: %v\n", c.App.Flags)
	fmt.Printf("Config: %#v\n", c.Context.Value(Config{}))
	fmt.Printf("Config RCPath: %s\n", c.Context.Value(Config{}).(Config).RCPath)
	fmt.Println("NYI")
}

func rcAction(c *cli.Context) error {
	arg := c.Args().First()
	var cloud Cloud
	clouds := c.Context.Value(contextKey("cloud")).([]Cloud)
	if arg == "" {
		cloud = rcselector(clouds)
	} else {
		cloud = GetCloud(arg, clouds)
		if cloud.Name == "" {
			return fmt.Errorf("Cloud %s not found.", arg)
		}
	}
	config := c.Context.Value(Config{}).(Config)
	config.SpawnEnv(cloud)
	return nil
}
func exportAction(c *cli.Context) error {
	debugNYI(c)
	return nil
}
func infoAction(c *cli.Context) error {
	debugNYI(c)
	return nil
}

func checkForNested() {
	env := os.Getenv(nestedEnvKey)
	if env == nestedEnvVal {
		fmt.Println("Exit current ossie session first.")
		os.Exit(0)
	}
}

func main() {

	checkForNested()

	app := &cli.App{
		Before: func(cctx *cli.Context) error {
			configfile := cctx.String("config")
			conf, _ := SetupConfig(configfile)
			cctx.Context = context.WithValue(cctx.Context, Config{}, conf)

			clouds := loadClouds(conf)
			if len(clouds) == 0 {
				return fmt.Errorf("No clouds found")
			}
			cctx.Context = context.WithValue(cctx.Context, contextKey("cloud"), clouds)
			return nil
		},
		Name:    "ossie",
		Usage:   "A powerful Tool to manage Openstack environments",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Usage: "Path to config `FILE`"},
		},
		Commands: []*cli.Command{
			{
				Name:      "rc",
				Usage:     "Spawn Shell with selected environment",
				Action:    rcAction,
				Args:      true,
				ArgsUsage: "[rc]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "region",
						Usage: "Sets region if available",
					},
				},
			},
			{
				Name:   "export",
				Usage:  "Export active environmen to stdout",
				Action: exportAction,
				Args:   false,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "mode",
						Usage:       "export mode: \"rc\" or \"yaml\"",
						DefaultText: "rc",
						Value:       "rc",
					},
				},
			},
			{
				Name:   "info",
				Usage:  "Shows active or selected environment",
				Action: infoAction,
				Args:   false,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
