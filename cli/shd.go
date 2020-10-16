package main

import (
	"fmt"
	"time"

	"github.com/yankghjh/selfhosted_store/cli/pipe"

	_ "github.com/yankghjh/selfhosted_store/cli/modules/app"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/docker-compose"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/icon"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/yacht"

	"github.com/spf13/viper"
)

var starttime time.Time
var cfg *viper.Viper

func init() {
	starttime = time.Now()
	cfg = viper.New()
	cfg.SetDefault("source", "apps")
	cfg.SetDefault("dist", "templates")
	cfg.SetDefault("icon.basepath", "https://raw.githubusercontent.com/yangkghjh/selfhosted_store/main/")
}

func main() {
	cfg.SetConfigName("config")
	cfg.SetConfigType("yml")
	cfg.AddConfigPath(".")

	err := cfg.ReadInConfig()
	if err != nil {
		fmt.Println("no conifg file found, use default config")
	}

	opt := pipe.Opition{
		SourcePath: cfg.GetString("source"),
		DistPath:   cfg.GetString("dist"),
		Sources:    []string{"app", "docker-compose", "icon"},
		Handlers:   []string{"yacht"},
		Config:     cfg,
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
