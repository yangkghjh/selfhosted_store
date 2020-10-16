package main

import (
	"fmt"
	"time"

	"github.com/yankghjh/selfhosted_store/cli/pipe"

	_ "github.com/yankghjh/selfhosted_store/cli/modules/app"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/docker-compose"
	_ "github.com/yankghjh/selfhosted_store/cli/modules/yacht"

	"github.com/spf13/viper"
)

var starttime time.Time

func init() {
	starttime = time.Now()
	viper.SetDefault("source", "apps")
	viper.SetDefault("dist", "templates")
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("no conifg file found, use default config")
	}

	opt := pipe.Opition{
		SourcePath: viper.GetString("source"),
		DistPath:   viper.GetString("dist"),
		Sources:    []string{"app", "docker-compose"},
		Handlers:   []string{"yacht"},
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
