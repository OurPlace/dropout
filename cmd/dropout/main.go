package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ourplace/dropout/internal/ping"
	"github.com/ourplace/dropout/internal/settings"
	"github.com/ourplace/dropout/internal/storage/postgres"
)

// Build Constants
const (
	Version   = "dev"
	GitCommit = "working"
)

func run(ctx context.Context) error {
	// Load our settings.
	config, err := settings.Load()
	if err != nil {
		return fmt.Errorf("Failed to load settings: %w", err)
	}

	// Open the database.
	db, err := postgres.Open(ctx, config.Database)
	if err != nil {
		return fmt.Errorf("Failed to open database: %w", err)
	}
	defer db.Close()

	// Perform our Ping.
	ping, err := ping.Perform(ctx, config.Ping)
	if err != nil {
		return fmt.Errorf("Failed to perform ping: %w", err)
	}

	// Exit if the ping was successful
	if ping {
		return nil
	}

	fmt.Println("Ping failed, reporting to the database.")
	// Report dropout.
	if err = db.ReportDropout(ctx); err != nil {
		return fmt.Errorf("Failed to store dropout: %w", err)
	}

	return nil
}

func main() {
	now := time.Now()
	fmt.Printf("Starting Dropout %s(%s)\nStarting at %v\n", Version, GitCommit, now.Format(time.RFC1123))
	// Create our timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Run and check for errors.
	if err := run(ctx); err != nil {
		cancel()

		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	finished := time.Now()
	timeItTook := finished.Sub(now)

	fmt.Printf("Finished at %v\nTime it took: %v\n", finished.Format(time.RFC1123), timeItTook)
}
