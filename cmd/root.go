package main

import (
	"log"
	"os"

	"github.com/fatihdumanli/onenote/internal/style"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long:                  "Take notes on your Onenote notebooks from terminal",
	Use:                   "nnote",
	DisableFlagsInUseLine: true,
	CompletionOptions:     cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("\n" + style.Error(err.Error()))
		os.Exit(1)
	}

}
