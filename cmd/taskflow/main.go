package main

import (
	"flag"
	"os"

	"github.com/egor-kardash/taskflow/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", os.Getenv("CONFIG_PATH"), "path to config file")
	flag.Parse()

	config.MustLoad(configPath)
}
