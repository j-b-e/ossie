package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/j-b-e/ossie/internal/config"
	"github.com/j-b-e/ossie/internal/model"
	"github.com/j-b-e/ossie/internal/shell"
	"github.com/urfave/cli/v3"
)

type contextKey string

func debugNYI(ctx context.Context, cmd *cli.Command) {
	fmt.Printf("Args: %#v\n", cmd.Args())
	fmt.Printf("Flagnames: %#v\n", cmd.FlagNames())
	fmt.Printf("LocalFlags: %#v\n", cmd.LocalFlagNames())
	fmt.Printf("Flags: %v\n", cmd.Flags)
	fmt.Printf("Config: %#v\n", ctx.Value(config.Config{}))
	fmt.Printf("Config RCPath: %s\n", ctx.Value(config.Config{}).(config.Config).RCPath)
	fmt.Println("NYI")
}

func rcAction(ctx context.Context, cmd *cli.Command) error {
	arg := cmd.Args().First()
	var cloud model.Cloud
	if arg == "" {
		cloud = selector(config.Global.Clouds)
	} else {
		cloud = config.Global.Clouds.Select(arg)
		if cloud.Name == "" {
			return fmt.Errorf("Cloud %s not found.", arg)
		}
	}

	shell.SpawnEnv(cloud)
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
	env := os.Getenv(config.NestedEnvKey)
	if env == config.NestedEnvVal {
		fmt.Println("Exit current ossie session first.")
		os.Exit(0)
	}
}

func Ossie(version string) {
	cmd := &cli.Command{
		Name:    "ossie",
		Usage:   "A powerful Tool to manage Openstack environments",
		Version: version,
		Before:  config.SetupConfig,
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
