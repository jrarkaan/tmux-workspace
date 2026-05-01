package config

type Config struct {
	Version    int                  `yaml:"version"`
	Defaults   Defaults             `yaml:"defaults"`
	Workspaces map[string]Workspace `yaml:"workspaces"`
}

type Defaults struct {
	Attach          bool `yaml:"attach"`
	CreateIfMissing bool `yaml:"create_if_missing"`
	BaseIndex       int  `yaml:"base_index"`
}

type Workspace struct {
	Root    string            `yaml:"root"`
	Env     map[string]string `yaml:"env"`
	Windows []Window          `yaml:"windows"`
}

type Window struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}
