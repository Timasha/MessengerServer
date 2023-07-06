package utils

import (
	"flag"
	"fmt"
)

type Flags struct {
	ConfigPath   string
	TemplatePath string
}

func ParceFlags() (flags Flags) {
	flag.StringVar(&flags.ConfigPath, "conf", "", "config path")
	flag.StringVar(&flags.TemplatePath, "templ", "", "template path")
	flag.Parse()
	fmt.Printf("Current config path: %v\n", flags.ConfigPath)
	return
}
