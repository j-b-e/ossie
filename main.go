package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var version = "dev"

const (
	nestedEnvKey = "__OSSIE_SPAWNED"
	nestedEnvVal = "righto"
)

type contextKey string

func debugNYI(ctx context.Context, cmd *cli.Command) {
	fmt.Printf("Args: %#v\n", cmd.Args())
	fmt.Printf("Flagnames: %#v\n", cmd.FlagNames())
	fmt.Printf("LocalFlags: %#v\n", cmd.LocalFlagNames())
	fmt.Printf("Flags: %v\n", cmd.Flags)
	fmt.Printf("Config: %#v\n", ctx.Value(Config{}))
	fmt.Printf("Config RCPath: %s\n", ctx.Value(Config{}).(Config).RCPath)
	fmt.Println("NYI")
}

func rcAction(ctx context.Context, cmd *cli.Command) error {
	arg := cmd.Args().First()
	var cloud Cloud
	if arg == "" {
		cloud = rcselector(gConf.clouds)
	} else {
		cloud = GetCloud(arg, gConf.clouds)
		if cloud.Name == "" {
			return fmt.Errorf("Cloud %s not found.", arg)
		}
	}

	SpawnEnv(cloud)
	return nil
}
func exportAction(ctx context.Context, cmd *cli.Command) error {
	debugNYI(ctx, cmd)
	return nil
}
func infoAction(ctx context.Context, cmd *cli.Command) error {
	debugNYI(ctx, cmd)
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
	cmd := &cli.Command{
		Name:    "ossie",
		Usage:   "A powerful Tool to manage Openstack environments",
		Version: version,
		Before:  SetupConfig,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config `FILE`",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "rc",
				Usage:     "Spawn Shell with selected environment",
				Action:    rcAction,
				ArgsUsage: "[rc]",

			},
			{
				Name:   "export",
				Usage:  "Export active environment to stdout",
				Action: exportAction,
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
			},
		},
	}

	checkForNested()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
