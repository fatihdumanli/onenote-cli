package main

import (
	"log"
	"os"

	"github.com/fatihdumanli/onenote"
	"github.com/fatihdumanli/onenote/internal/style"
	"github.com/fatihdumanli/onenote/internal/survey"
	"github.com/fatihdumanli/onenote/pkg/msftgraph"
	errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: "Take notes on your Onenote notebooks from terminal",
	RunE: func(c *cobra.Command, args []string) error {
		var code, err = startNoteSurvey()
		os.Exit(code)
		return err
	},
	Use:                   "cnote",
	DisableFlagsInUseLine: true,
}

//The function gets executed once the application starts without any commands/arguments.
func startNoteSurvey() (_ int, err error) {

	defer func() {
		if err != nil {
			err = errors.Wrap(err, "in startNoteSurvey\n")
		}
	}()

	noteContent, err := survey.AskNoteContent()
	if err != nil {
		return 1, err
	}

	notebooks, err := onenote.GetNotebooks()
	if err != nil {
		return 2, err
	}

	n, err := survey.AskNotebook(notebooks)
	if err != nil {
		return 1, err
	}
	sections, err := onenote.GetSections(n)
	if err != nil {
		return 3, err
	}
	section, err := survey.AskSection(n, sections)
	if err != nil {
		return 4, err
	}

	title, err := survey.AskTitle()
	if err != nil {
		return 4, err
	}

	//Saving the note to the section
	_, err = onenote.SaveNotePage(*msftgraph.NewNotePage(section, title, noteContent), false)
	if err != nil {
		return 1, err
	}

	return 0, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(style.Error(err.Error()))
		os.Exit(1)
	}

}
