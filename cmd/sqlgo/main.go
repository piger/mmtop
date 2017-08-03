package main

import (
	"flag"
	"fmt"
	"github.com/piger/sqlgo"
	"os"
)

var configFilename = flag.String("config", "config.ini", "Configuration filename")

func main() {
	flag.Parse()

	configs, err := sqlgo.ReadConfig(*configFilename)
	if err != nil {
		fmt.Printf("Error in configuration: %s\n", err)
		os.Exit(1)
	}

	sqlgo.RunClient(configs)
}
