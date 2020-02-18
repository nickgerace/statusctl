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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"

	git "gopkg.in/src-d/go-git.v4"
	yaml "gopkg.in/yaml.v3"
)

// Config struct for yaml files. Only look for collections of Git repositories and individual Git
// repositories. The "omitempty" option allows for the objects to be empty.
type Config struct {
	Collections  []string `yaml:"collections,omitempty"`
	Repositories []string `yaml:"repositories,omitempty"`
}

// Help message displaying whenever the flag.Usage() function is overridden.
var help string = (`
Statusctl
https://github.com/nickgerace/statusctl

description:
  CLI tool to keep track of your Git repositories. It leverages the configuration YAML file
  ($HOME/.config/statusctl/config.yaml) in order to know which repositories to target.

config file:
  The file is split into *collections* and *repositories* containing arrays of paths. Each entry in
  *collections* is a path to a collection of Git repositories. Each entry in *repositories* is a path
  to an individual Git repository.

usage:
`)

// Open the file and load the config into the struct for usage. We need a pointer to the config
// object to store the unmarhsaled data. The config file can only be in one location. First, we
// load the config file contents into bytestream. Then, we unmarshall the raw data into the
// config struct.
func (config *Config) load() {
	data, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".config", "statusctl", "config.yaml"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Checks if the config file exists and returns the result. We use the directory path as well as
// the config file path in order to create the file if it does not exist, create the nested
// directories if they do not exist, and exit afterwards.
func createConfigIfNotExist() {
	configDirectory := path.Join(os.Getenv("HOME"), ".config", "statusctl")
	if _, err := os.Stat(path.Join(configDirectory, "config.yaml")); err != nil {
		err = os.MkdirAll(configDirectory, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, err = os.Create(path.Join(configDirectory, "config.yaml"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}

// List all items in the current config. This function starts by creating the config file (and
// exit) if it does not exist. Then, it loads the YAML file into a config struct. Finally, it
// prints all config contents to STDOUT.
func listAction() {
	createConfigIfNotExist()
	config := Config{}
	config.load()
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

// With a slice of repositories, print the status of each repository using the go-git library.
// Using an array of STDOUT options, iterate through the slice of repositories, and open the
// repository with a given path (transformed to absolute path). This will not only fail if the
// path does not exist, but also if the path is not a valid (Git) repository. The "IsClean"
// function requires creating a "worktree" object, and then a "status" object. If either function
// fails, the error is unknown since the previous function should have caught most known errors.
// Finally, we will utilize the result (bool) to determine whether or not the repository is clean.
func printStatus(repos []string) {
	var waitGroup sync.WaitGroup
	results := [4]string{
		"CLEAN    ",
		"UNCLEAN  ",
		"ERROR    ",
		"UNKNOWN  ",
	}
	for _, repoPath := range repos {
		waitGroup.Add(1)
		go func(repoPath string) {
			defer waitGroup.Done()
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
			if GitStatus.IsClean() {
				fmt.Printf("  %s%s\n", results[0], repoPath)
			} else {
				fmt.Printf("  %s%s\n", results[1], repoPath)
			}
		}(repoPath)
	}
	waitGroup.Wait()
}

// This is the primary function for the run command. First, this function creates the config file
// (and exits) if it does not exist. Then, it loads the YAML file into the config struct. We
// iterate through all collections and get their paths; subsequently getting all subdirectories
// for each collection. For each subdirectory in a collection, this function gets the absolute
// path of the subdirectory and adds that path to a temporary slice. After this is done for every
// subdirectory, this function feeds the completed slice into "printStatus". Finally, all the
// individual repositories are fed directly into "printStatus".
func runAction() {
	createConfigIfNotExist()
	config := Config{}
	config.load()
	fmt.Printf("\ncollections:\n")
	for _, collection := range config.Collections {
		collectionSubDirs, err := ioutil.ReadDir(collection)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		subDirPaths := []string{}
		for _, subDir := range collectionSubDirs {
			subDirPath := path.Join(collection, subDir.Name())
			subDirPath, err = filepath.Abs(subDirPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			subDirPaths = append(subDirPaths, subDirPath)
		}
		printStatus(subDirPaths)
	}
	fmt.Printf("\nrepositories:\n")
	printStatus(config.Repositories)
	fmt.Printf("\n")
}

// Execute the main flag commands.
func main() {
	var run = flag.Bool("r", false, "run the main program: checking Git status for all collections and repositories")
	var list = flag.Bool("l", false, "list the contents of the configuration file ($HOME/.config/statusctl/config.yaml)")
	flag.Usage = func() {
		fmt.Printf(help)
		flag.PrintDefaults()
		fmt.Printf("\n")
		os.Exit(0)
	}
	flag.Parse()
	if (*run == false) && (*list == false) {
		flag.Usage()
	}
	if (*run == true) && (*list == true) {
		fmt.Printf("Cannot use more than one flag. Printing help message...\n")
		flag.Usage()
	}
	if *run == true {
		runAction()
	}
	if *list == true {
		listAction()
	}
}
