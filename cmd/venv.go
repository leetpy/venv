package main

import (
	"flag"
	"fmt"
	"github.com/leetpy/venv/internal"
	"github.com/leetpy/venv/manager"
	"os"
)

var (
	flagSet = flag.NewFlagSet("venv", flag.ExitOnError)

	config = flagSet.String("config", "", "path to config file")
)

func startSite(cfg *manager.Config) error {
	m := manager.GetManager(cfg)
	m.Run()
	return nil
}

func main() {
	flagSet.Parse(os.Args[1:])
	if *config == "" {
		*config = "/etc/venv/venv.conf"
	}

	if internal.Exist(*config) {
		cfg := manager.LoadConfig(*config)
		fmt.Println(*config)
		fmt.Println(cfg.Server.Port)
		fmt.Println(cfg.Debug)
		startSite(cfg)
	} else {
		os.Exit(1)
	}

}
