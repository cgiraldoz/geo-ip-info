package main

import (
	"github.com/cgiraldoz/geo-ip-info/cmd/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
