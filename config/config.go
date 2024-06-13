package config

// Configurations exported
type Configurations struct {
	Processes []ProcessConfigurations
}

// OutputConfigurations exported
type OutputConfigurations struct {
	Type   string
	Format string
	Path   string
}

// ProcessConfigurations exported
type ProcessConfigurations struct {
	Name    string
	Module  string
	Query   string
	Outputs []OutputConfigurations
}
