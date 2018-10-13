package main

import (
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pe",
		Short: "PE file resource extractor",
	}
	return cmd
}
