package main

import (
	"fmt"
	"os"

	"github.com/okaryo/dlog/cmd"
	"github.com/okaryo/dlog/internal/service"
	"github.com/okaryo/dlog/internal/storage"
)

func main() {
	store, err := storage.NewDefaultStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	svc := service.New(store)
	rootCmd := cmd.NewRootCmd(svc, os.Stdout, os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
