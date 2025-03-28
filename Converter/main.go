package main

import (
	"v1/currency"
	"v1/ui"
)

func main() {
	service := currency.NewService()
	cli := ui.NewCLI(*service)
	cli.RUN()
}
