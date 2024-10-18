package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func debugNYI(c *cli.Context) {
	fmt.Printf("%#v\n", c.Args())
	fmt.Printf("%#v\n", c.Context.Value(Config{}))
	fmt.Printf("%s\n", c.Context.Value(Config{}).(Config).RCPath)
	fmt.Println("NYI")
}

func rcAction(c *cli.Context) error {
	arg := c.Args().First()
	if arg == "" {
		fmt.Println("Select Cloud env")
		arg = rcselector(c.Context.Value(Clouds{}).([]Clouds))
	}
	SpawnEnv(arg)
	return nil
}
func regionAction(c *cli.Context) error {
	debugNYI(c)
	return nil
}
func editAction(c *cli.Context) error {
	debugNYI(c)
	return nil
}
func apiverAction(c *cli.Context) error {
	debugNYI(c)
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
func createAction(c *cli.Context) error {
	debugNYI(c)
	return nil
}

func main() {

	app := &cli.App{
		Before: func(cctx *cli.Context) error {
			configfile := cctx.String("config")
			conf, _ := SetupConfig(configfile)
			cctx.Context = context.WithValue(cctx.Context, Config{}, conf)

			clouds := loadCloudsYaml()
			if len(clouds) == 0 {
				fmt.Println("No clouds found.")
				return fmt.Errorf("No clouds found")
			}
			cctx.Context = context.WithValue(cctx.Context, Clouds{}, clouds)
			return nil
		},
		Name:    "ossie",
		Usage:   "A powerful Tool to manage Openstack Contexts",
		Version: "v1.0-dev",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Usage: "Path to config `FILE`"},
		},
		Commands: []*cli.Command{
			{
				Name:      "rc",
				Usage:     "set env to rc",
				Action:    rcAction,
				Args:      true,
				ArgsUsage: "[rc]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "region",
						Usage: "sets region if available",
					},
				},
			},
			{
				Name:   "regions",
				Usage:  "list regions",
				Action: regionAction,
				Args:   false,
			},
			{
				Name:      "edit",
				Usage:     "edit rc",
				Action:    editAction,
				Args:      true,
				ArgsUsage: "[rc]",
			},
			{
				Name:   "api-version",
				Usage:  "set api-version",
				Action: apiverAction,
				Args:   false,
			},
			{
				Name:   "export",
				Usage:  "export current active rc",
				Action: exportAction,
				Args:   false,
			},
			{
				Name:   "info",
				Usage:  "shows current or selected rc",
				Action: infoAction,
				Args:   false,
			},
			{
				Name:   "create",
				Usage:  "menu driven creation of rc from scratch or based on rc",
				Action: createAction,
				Args:   false,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
