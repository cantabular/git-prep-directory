package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sensiblecodeio/git-prep-directory"

	"github.com/urfave/cli/v2"
)

// CloneTimeout specifies the duration allowed for each individual `git clone`
// call (main repository mirroring or git submodule initialization) before
// cancelling the operation.
const CloneTimeout = 1 * time.Hour

func init() {
	log.SetFlags(0)
}

func main() {
	app := cli.NewApp()
	app.Name = "git-prep-directory"
	app.Version = "0.6.0"
	app.Usage = "Build tools friendly way of repeatedly cloning a git\n" +
		"   repository using a submodule cache and setting file timestamps to commit times."

	app.Action = actionMain

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "url, u",
			Usage: "URL to clone",
		},
		&cli.StringFlag{
			Name:  "ref, r",
			Usage: "ref to checkout",
		},
		&cli.StringFlag{
			Name:  "destination, d",
			Usage: "destination dir",
			Value: "./src",
		},
		&cli.DurationFlag{
			Name:    "timeout, t",
			Usage:   "clone timeout",
			Value:   CloneTimeout,
			EnvVars: []string{"GIT_PREP_DIR_TIMEOUT"},
		},
	}

	app.RunAndExitOnError()
}

func actionMain(c *cli.Context) error {
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
