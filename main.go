package main

import (
	"github.com/guionardo/gs-bot/cmd"
)

func main() {
	svc := cmd.GetBot()

	svc.Start()
}
