package main

import (
	"github.com/spf13/cobra"

	"github.com/yankghjh/selfhosted_store/cli/cmd/generate"
)

var (
	command = &cobra.Command{
		Use:   "shctl",
		Short: "Shctl is a selfhosted store generater",
	}
)

func init() {
	command.AddCommand(generate.Command)
}

func main() {
	command.Execute()
}
