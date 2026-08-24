package config

// Options configures the coldchain server process.
type Options struct {
	Addr    string
	DataDir string
}

// Default returns the development defaults for a local instance.
func Default() Options {
	return Options{Addr: ":8080", DataDir: "data"}
}

// WithAddr returns a copy of the options using the given listen address.
func (o Options) WithAddr(addr string) Options {
	o.Addr = addr
	return o
}

// WithDataDir returns a copy of the options using the given data directory.
func (o Options) WithDataDir(dir string) Options {
	o.DataDir = dir
	return o
}
