package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/philandstuff/dhall-golang/v5"
	"github.com/urfave/cli/v2"
)

type Script struct {
	Name    string `dhall:"name"`
	Command string `dhall:"command"`
	Context string `dhall:"context"`
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get current working directory")
	}

	configDir, err := findContainingFolderOfFileInWdOrParents(wd, "turbine.dhall")
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:  "turbine",
		Usage: "manage your project",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "run a script defined in turbine.dhall",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "The name of the script to run",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					config := readConfig(configDir)

					var script Script
					found := false

					for _, theScript := range config.Scripts {
						if theScript.Name == c.String("name") {
							found = true
							script = theScript
						}
					}

					if !found {
						log.Fatal("The specified script is not defined in turbine.dhall")
					}

					parts := strings.Fields(script.Command)
					head := parts[0]
					parts = parts[1:]

					cmd := exec.Command(head, parts...)
					cmd.Dir = filepath.Join(configDir, script.Context)
					stdout, err := cmd.Output()

					if err != nil {
						log.Fatal(err)
					}

					print(string(stdout))

					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Scripts []Script `dhall:"scripts"`
}

func readConfig(configDir string) Config {
	var config Config

	bytes, err := ioutil.ReadFile(filepath.Join(configDir, "turbine.dhall"))
	if err != nil {
		panic(err)
	}
	err = dhall.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func findContainingFolderOfFileInWdOrParents(path string, filename string) (string, error) {

	path, err := filepath.Abs(path)

	if err != nil {
		return "", fmt.Errorf("Couldn't find %s in current directory or any of it's parents", filename)
	}

	file := filepath.Join(path, filename)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		if path == "/" {
			return "", fmt.Errorf("Couldn't find %s in current directory or any of it's parents", filename)
		}
		return findContainingFolderOfFileInWdOrParents(filepath.Join(path, ".."), filename)
	}

	return path, nil

}
