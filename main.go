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
	"os"
)

// This is the base command, "statusctl", when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "statusctl",
	Short: "A brief description of your application",
	Long:  "Multi-line description here",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
