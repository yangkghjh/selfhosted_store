package main

import (
	"fmt"
	"os"
	"time"

	"github.com/yankghjh/selfhosted_store/cli/pipe"

	_ "github.com/yankghjh/selfhosted_store/cli/modules/app"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/docker-compose"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/icon"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/unraid"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/yacht"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	starttime time.Time
	cfg       *viper.Viper
	rootCmd   = &cobra.Command{
		Use:   "shctl",
		Short: "Shctl is a selfhosted store generator",
		Run:   execute,
	}
)

func init() {
	starttime = time.Now()
	cfg = viper.New()

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yml", "Config file (default is config.json)")

	cfg.SetDefault("source", "apps")
	cfg.SetDefault("dist", "templates")
	cfg.SetDefault("icon.basepath", "https://yangkghjh.github.io/selfhosted_store/")
	// cfg.SetDefault("sources", []string{"app", "docker-compose", "icon"})
	// cfg.SetDefault("handlers", []string{"yacht"})
	cfg.SetDefault("skipSourceFile", false)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute(cmd *cobra.Command, args []string) {
	cfg.SetConfigFile(cfgFile)

	err := cfg.ReadInConfig()
	if err != nil {
		fmt.Println("no conifg file found, use default config")
	}

	opt := pipe.Opition{
		SourcePath:     cfg.GetString("source"),
		DistPath:       cfg.GetString("dist"),
		Sources:        cfg.GetStringSlice("sources"),
		Handlers:       cfg.GetStringSlice("handlers"),
		SkipSourceFile: cfg.GetBool("skipSourceFile"),
		Config:         cfg,
	}

	p, err := pipe.NewPipe(opt)
	if err != nil {
		fmt.Println("create pipe error: ", err)
		return
	}

	err = p.Run()
	if err != nil {
		fmt.Println("run pipe error: ", err)
		return
	}

	fmt.Printf("parsed %d apps in %s\n", p.Count, time.Now().Sub(starttime))
}
