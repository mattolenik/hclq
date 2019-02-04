package config

// Configuration settings
type Configuration struct {
	UseRawOutput bool
	InputFile    string
	OutputFile   string
}

// Settings app configuration
var Settings *Configuration

func init() {
	Settings = &Configuration{}
}
