package main

import "os"

func main() {
	rootCmd := newRootCommand()
	rootCmd.AddCommand(
		newExtractCommand(),
		newExtractAllCommand(),
		newListCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
