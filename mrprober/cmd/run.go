package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"mrprober/engine"
)

func newRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "One-shot execution of the rules.",
		Args:  cobra.ExactArgs(0),
		Run:   runCommandFunc,
	}

	return cmd
}

func runCommandFunc(cmd *cobra.Command, args []string) {

	// Execute and print results
	for r := range engine.OneShotRun() {
		log.Println(r)
	}
}
