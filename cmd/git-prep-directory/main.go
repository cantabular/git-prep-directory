package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sensiblecodeio/git-prep-directory"

	"github.com/urfave/cli/v3"
)

// CloneTimeout specifies the duration allowed for each individual `git clone`
// call (main repository mirroring or git submodule initialization) before
// cancelling the operation.
const CloneTimeout = 1 * time.Hour

func init() {
	log.SetFlags(0)
}

func main() {
	cmd := &cli.Command{
		Name:    "git-prep-directory",
		Version: "0.7.0",
		Usage: "Build tools friendly way of repeatedly cloning a git\n" +
			"   repository using a submodule cache and setting file timestamps to commit times.",
		Action: actionMain,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "url",
				Aliases: []string{"u"},
				Usage: "URL to clone",
			},
			&cli.StringFlag{
				Name:  "ref, r",
				Aliases: []string{"r"},
				Usage: "ref to checkout",
			},
			&cli.StringFlag{
				Name:  "destination",
				Aliases: []string{"d"},
				Usage: "destination dir",
				Value: "./src",
			},
			&cli.DurationFlag{
				Name:    "timeout",
				Aliases: []string{"t"},
				Usage:   "clone timeout",
				Value:   CloneTimeout,
				Sources: cli.EnvVars("GIT_PREP_DIR_TIMEOUT"),
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func actionMain(_ context.Context, c *cli.Command) error {
	if !c.IsSet("url") || !c.IsSet("ref") {
		log.Fatalln("Error: --url and --ref required")
	}

	where, err := git.PrepBuildDirectory(
		c.String("destination"),
		c.String("url"),
		c.String("ref"),
		c.Duration("timeout"),
		os.Stderr)
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Printf("Checked out %v at %v", where.Name, where.Dir)
	fmt.Println(where.Dir)

	return nil
}
