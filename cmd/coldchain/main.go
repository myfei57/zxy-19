package main

import (
	"flag"
	"log"

	"coldchain/internal/config"
	"coldchain/internal/console"
	"coldchain/internal/store"
)

func main() {
	addr := flag.String("addr", "", "HTTP listen address")
	dataDir := flag.String("data", "", "data directory")
	seed := flag.Bool("seed", false, "seed a demo cold-chain dataset")
	flag.Parse()

	opts := config.Default()
	if *addr != "" {
		opts = opts.WithAddr(*addr)
	}
	if *dataDir != "" {
		opts = opts.WithDataDir(*dataDir)
	}
	if err := opts.Validate(); err != nil {
		log.Fatalf("coldchain: %v", err)
	}

	st, err := store.Open(opts.DataDir)
	if err != nil {
		log.Fatalf("coldchain: %v", err)
	}

	services := buildServices(st)
	if *seed {
		if err := seedDemo(services); err != nil {
			log.Fatalf("coldchain: seed: %v", err)
		}
	}

	server := console.New(opts, st, services)
	log.Printf("coldchain listening on %s (data=%s)", opts.Addr, opts.DataDir)
	if err := server.Start(); err != nil {
		log.Fatalf("coldchain: server: %v", err)
	}
}
