package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/folbricht/pefile"
	"github.com/spf13/cobra"
)

func newExtractAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extract-resources <file> <dir>",
		Short: "Extract all resources from a PE file",
		Long: `Extract all resources from a PE file into the given
output directory. The directory structure in the PE file is
preserved in the output directory.`,
		Example: `  pe extract-resources file.exe /tmp/resources`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExtractAll(args)
		},
		SilenceUsage: true,
	}
	return cmd
}

func runExtractAll(args []string) error {
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
		dst := filepath.Join(args[1], r.Name)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		if err := ioutil.WriteFile(dst, r.Data, 0644); err != nil {
			return err
		}
	}
	return nil
}
