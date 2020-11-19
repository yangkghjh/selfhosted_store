package convert

import (
	"fmt"
	"io/ioutil"

	"github.com/yankghjh/selfhosted_store/cli/converter"

	"github.com/spf13/cobra"
)

var (
	// Command for convert
	Command = &cobra.Command{
		Use:   "convert",
		Short: "convert format between compose files and templates",
		Run:   run,
	}
	input string
	from  string
	to    string
)

func init() {
	Command.Flags().StringVarP(&input, "input", "i", "docker-compose.yml", "input file")
	Command.Flags().StringVarP(&from, "from", "f", "docker-compose", "source format")
	Command.Flags().StringVarP(&to, "to", "t", "docker-compose", "target format")
}

func run(cmd *cobra.Command, args []string) {
	payload, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Printf("read input file %s error: %s\n", input, err.Error())
		return
	}

	output, err := converter.Convert(from, to, payload)
	if err != nil {
		fmt.Printf("convert error: %s\n", err.Error())
		return
	}

	fmt.Println(string(output))
}
