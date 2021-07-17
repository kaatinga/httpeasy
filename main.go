package httpeasy

// NewConfig creates new config model that later can be processed
// by settings package or updated manually.
func NewConfig() *Config {
	return new(Config)
}
