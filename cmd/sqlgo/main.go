package main

import (
	"fmt"
	"github.com/piger/sqlgo"
	"os"
)

func main() {
	configs, err := sqlgo.ReadConfig("config.ini")
	if err != nil {
		fmt.Printf("Error in configuration: %s\n", err)
		os.Exit(1)
	}

	for _, config := range configs {
		fmt.Printf("mysql://%s:%s@%s:%d\n", config.Username, config.Password, config.Address, config.Port)
	}

	sqlgo.RunClient(configs)
}
