package config

import "fmt"

// Validate checks that the options are usable by the server.
func (o Options) Validate() error {
	if o.Addr == "" {
		return fmt.Errorf("config: listen address is empty")
	}
	if o.DataDir == "" {
		return fmt.Errorf("config: data directory is empty")
	}
	return nil
}
