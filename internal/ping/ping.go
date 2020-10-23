package ping

import (
	"context"
	"fmt"

	// We're going to use go-ping/ping to save time.
	"github.com/go-ping/ping"
	"github.com/ourplace/dropout/internal/settings"
)

// Perform performs a ping and reports true if it was successful.
func Perform(ctx context.Context, config settings.Ping) (bool, error) {
	// Loop over each Location
	for _, location := range config.Locations {
		// Attempt to ping the location.
		successful, err := pingIP(ctx, config, location)
		if err != nil {
			return false, err
		}
		// If we found it, we don't need to try any further.
		if successful {
			return true, nil
		}
		// If it wasn't successful, continue the loop.
	}

	return false, nil
}

// pingIP will attempt to Ping the address. Not being able to find the route to
// host is not considered an error, but a failure to Ping.
func pingIP(ctx context.Context, config settings.Ping, location string) (successful bool, err error) {
	pinger, err := ping.NewPinger(location)
	if err != nil {
		return false, err
	}

	pinger.Count = 1
	pinger.Timeout = config.Timeout
	pinger.SetPrivileged(config.Privileged)

	if err := pinger.Run(); err != nil {
		return false, err
	}

	stats := pinger.Statistics()

	fmt.Printf("Ping took %s\n", stats.MaxRtt)

	if stats.PacketsRecv != 1 {
		return false, nil
	}

	return true, nil
}
