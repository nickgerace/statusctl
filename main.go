/*
 * STATUSCTL
 * created by: Nick Gerace
 *
 * MIT License, Copyright (c) Nick Gerace
 * See "LICENSE" file for more information.
 *
 * Please find license and further information via the link below.
 * https://github.com/nickgerace/statusctl
 */

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path"
)

// Config struct for yaml files. Only look for collections of git repositories and individual git
// repositories. The "omitempty" option enables the results to be in a line-by-line list rather
// than an inline array. This is done for readability.
type Config struct {
	Collections  []string `yaml:"collections,omitempty"`
	Repositories []string `yaml:"repositories,omitempty"`
}

// This is the base command, "statusctl", when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "statusctl",
	Short: "A brief description of your application",
	Long:  "Multi-line description here",
}

// List all targets to check status against.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  "Multi-line description here",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Yo")
	},
}

// Run againsts all targets in the config YAML file.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long:  "Multi-line description here",
	Run: func(cmd *cobra.Command, args []string) {
		config := Config{}
		config.load()
	},
}

// Open the file and load the config in to the struct for usage. We need a pointer to the config
// object to store the unmarhsaled data. The config file can only be in one location.
func (config *Config) load() {

	// Load the config file contents into bytestream.
	data, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".config", "statusctl", "config.yaml"))
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshall the raw string into the config struct.
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
}

// Checks if the config file exists and returns the result.
func checkConfigExists() bool {
	if _, err := os.Stat(path.Join(os.Getenv("HOME"), ".config", "statusctl", "config.yaml")); err != nil {
		return false
	}
	return true
}

// With a slice of repositories, print the status of each repository using the go-git library.
func printStatus(repos []string) {

	// Types of results and formatted output.
	results := [4]string{
		"CLEAN    ",
		"UNCLEAN  ",
		"ERROR    ",
		"UNKNOWN  ",
	}

	// Iterate through the slice of repositories.
	for _, repoPath := range repos {

		// Open the repository with a given path. This not only fails if the path does not exist,
		// but also if the path is not a valid (git) repository. Before that, we will use the
		// absolute path.
		repoPath, err = path.filepath.Abs(repoPath)
		if err != nil {
			log.Fatal(err)
		}
		GitRepo, err := git.PlainOpen(repoPath)
		if err != nil {
			fmt.Printf("  %s%s\n", results[2], repoPath)
			return
		}

		// The "IsClean" function requires creating a "worktree" object, and then a "status"
		// object. If either function fails, the error is unknown since the previous function
		// should have caught most known errors.
		GitWorktree, err := GitRepo.Worktree()
		if err != nil {
			fmt.Printf("  %s%s: %v\n", results[3], repoPath, err)
			return
		}
		GitStatus, err := GitWorktree.Status()
		if err != nil {
			fmt.Printf("  %s%s: %v\n", results[3], repoPath, err)
			return
		}

		// Utilize the result (bool) to provide whether or not the repository is clean.
		if GitStatus.IsClean() {
			fmt.Printf("  %s%s\n", results[0], repoPath)
		} else {
			fmt.Printf("  %s%s\n", results[1], repoPath)
		}
	}
}

// Setup all subcommands.
func init() {
	rootCmd.AddCommand(listCmd)
}

// Execute the root Cobra command.
func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
