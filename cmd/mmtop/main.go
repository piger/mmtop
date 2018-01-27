// Main entry point for mmtop.
package main

import (
	"flag"
	"fmt"
	"github.com/piger/mmtop"
	"os"
)

var configFilename = flag.String("config", "mmtop.ini", "Read config from FILE")

func main() {
	flag.Parse()

	configs, err := mmtop.ReadConfig(*configFilename)
	if err != nil {
		fmt.Printf("ERROR in configuration file '%s'\n%s\n", *configFilename, err)
		os.Exit(1)
	}

	mmtop.RunClient(configs)
}
