package generate

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yankghjh/selfhosted_store/cli/project"

	// modules
	_ "github.com/yankghjh/selfhosted_store/cli/modules/app"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/unraid"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/yacht"
)

var (
	// Command for generate
	Command = &cobra.Command{
		Use:   "generate",
		Short: "generate site and templates by configure file",
		Run:   run,
	}
	cfgFile string
	cfg     *viper.Viper
)

func init() {
	cfg = viper.New()
	Command.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yml", "Config file (default is config.json)")
}

func run(cmd *cobra.Command, args []string) {
	starttime := time.Now()
	cfg.SetConfigFile(cfgFile)

	err := cfg.ReadInConfig()
	if err != nil {
		fmt.Println("no conifg file found, use default config")
		return
	}

	p := project.NewProject(cfg)

	loaders := cfg.GetStringMap("loaders")
	if len(loaders) == 0 {
		fmt.Println("no loader found in config file")
		return
	}
	for name := range loaders {
		err := p.AddLoader(name)
		if err != nil {
			fmt.Printf("parse loader error: %s\n", err)
			return
		}
	}

	generaters := cfg.GetStringMap("generaters")
	for name := range generaters {
		err := p.AddGenerater(name)
		if err != nil {
			fmt.Printf("parse generater error: %s\n", err)
			return
		}
	}

	err = p.Run()
	if err != nil {
		fmt.Printf("generate error: %s\n", err)
		return
	}

	fmt.Printf("parsed %d apps in %s\n", len(p.Apps), time.Now().Sub(starttime))
}
