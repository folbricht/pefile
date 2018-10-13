package main

import (
	"fmt"

	"github.com/folbricht/pefile"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-resources <file>",
		Short:   "List resources in a PE file",
		Example: `  pe list-resources file.exe`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(args)
		},
		SilenceUsage: true,
	}
	return cmd
}

func runList(args []string) error {
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
		fmt.Println(r.Name)
	}
	return nil
}
