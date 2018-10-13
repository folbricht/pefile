package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/folbricht/pefile"
	"github.com/spf13/cobra"
)

func newExtractCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-resource <PEfile> <resourcename> [<outputfile>]",
		Short: "Extract a single resource from a PE file",
		Long: `Extract a single resource from a PE by name or ID. If
no output file name is givev, it'll be written to STDOUT.
Use 'list-resources' to get a list of resources available for
extraction.`,
		Example: `  pe extract-resourse file.exe 10/SOMERESOURCE/1033 someresource.bin`,
		Args:    cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExtract(args)
		},
		SilenceUsage: true,
	}
	return cmd
}

func runExtract(args []string) error {
	f, err := pefile.Open(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	resources, err := f.GetResources()
	if err != nil {
		return err
	}
	for _, r := range resources {
		if r.Name == args[1] {
			if len(args) > 2 { // Got an output file? Write to that
				return ioutil.WriteFile(args[2], r.Data, 0644)
			}
			// Write to stderr
			_, err := os.Stdout.Write(r.Data)
			return err
		}
	}
	return fmt.Errorf("resource '%s' not found in %s", args[1], args[0])
}
