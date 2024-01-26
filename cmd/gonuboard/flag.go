package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dukhyungkim/gonuboard/version"
)

var (
	FlagAddr    string
	FlagVersion bool
	FlagHelp    bool
)

func parseFlags() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&FlagAddr, "addr", "127.0.0.1:8080", "address")
	flag.BoolVar(&FlagVersion, "version", false, "print version")
	flag.BoolVar(&FlagHelp, "help", false, "print help message")
	flag.Parse()
}

func printVersion() {
	fmt.Println(version.Version)
}
