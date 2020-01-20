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
	"os"
	"path"
	"path/filepath"
)

// Config struct for yaml files. Only look for collections of git repositories and individual git
// repositories. The "omitempty" option enables the results to be in a line-by-line list rather
// than an inline array. This is done for readability.
type Config struct {
	Collections  []string `yaml:"collections,omitempty"`
	Repositories []string `yaml:"repositories,omitempty"`
}

// This is the base command, "statusctnl", when called without any subcommands.
// TODO: command descriptions
var rootCmd = &cobra.Command{
	Use:   "statusctl",
	Short: "Status ctl keeps track of ",
	Long:  "Multi-line description here",
}

// List all targets to check status against.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long:  "Multi-line description here",
	Run: func(cmd *cobra.Command, args []string) {
		listAction()
	},
}

// Run againsts all targets in the config YAML file.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long:  "Multi-line description here",
	Run: func(cmd *cobra.Command, args []string) {
		runAction()
	},
}

// Open the file and load the config in to the struct for usage. We need a pointer to the config
// object to store the unmarhsaled data. The config file can only be in one location.
func (config *Config) load() {

	// Load the config file contents into bytestream.
	data, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".config", "statusctl", "config.yaml"))
	handle(err)

	// Unmarshall the raw string into the config struct.
	err = yaml.Unmarshal(data, &config)
	handle(err)
}

// Handles errors efficiently in one line.
func handle(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Checks if the config file exists and returns the result.
func createConfigIfNotExist() {

	// Declare the config directory path.
	configDirectory := path.Join(os.Getenv("HOME"), ".config", "statusctl")

	// If the config file does not exist, create the file.
	if _, err := os.Stat(path.Join(configDirectory, "config.yaml")); err != nil {

		// First, create the nested directories to the file.
		err = os.MkdirAll(configDirectory, os.ModePerm)
		handle(err)

		// Second, create the file and exit.
		_, err = os.Create(path.Join(configDirectory, "config.yaml"))
		handle(err)
		os.Exit(0)
	}
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
		repoPath, err := filepath.Abs(repoPath)
		if err != nil {
			fmt.Printf("  %s%s: %v\n", results[3], repoPath, err)
			return
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

// This is the primary function for the run command. It iterates through all the collections and
// then gets the absolute paths of all the sub directories in each collection. Those paths are
// fed into "printStatus" as a slice. Then, the individual repositories are fed directory into
// "printStatus".
func runAction() {

	// Create config file (and exit) if it does not exist.
	createConfigIfNotExist()

	// Load YAML file into config struct.
	config := Config{}
	config.load()

	// Iterates through all collections and get their paths.
	fmt.Printf("\ncollections:\n")
	for _, collection := range config.Collections {

		// Get all subdirectories in the collection.
		collectionSubDirs, err := ioutil.ReadDir(collection)
		handle(err)

		// Iterate through the subdirectories in the collection.
		for count, subDir := range collectionSubDirs {

			// Get the absolute path of the sub directory.
			subDir, err = filepath.Abs(path.Join(collection, subDir.Name()))
			handle(err)

			// Replace the subdirectory name with the absolute path of the subdirectory.
			collectionSubDirs[count] = subDir
		}

		// Now that we have the absolute paths, feed them into "printStatus".
		printStatus(collectionSubDirs)
	}

	// Feed all individual repositories into "printStatus".
	fmt.Printf("\nrepositories:\n")
	printStatus(config.Repositories)
	fmt.Printf("\n")
}

// List all items in the current config.
func listAction() {

	// Create config file (and exit) if it does not exist.
	createConfigIfNotExist()

	// Load YAML file into config struct.
	config := Config{}
	config.load()

	// Perform the listing of the YAML contents.
	fmt.Printf("\ncollections:\n")
	for _, collection := range config.Collections {
		fmt.Printf("  %s\n", collection)
	}
	fmt.Printf("\nrepositories:\n")
	for _, repository := range config.Repositories {
		fmt.Printf("  %s\n", repository)
	}
	fmt.Printf("\n")
}

// Setup all subcommands.
func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(runCmd)
}

// Execute the root Cobra command.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
