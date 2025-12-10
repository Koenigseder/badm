package packages

type ConfigFile struct {
	Packages Package           `yaml:"packages"`
	Scripts  map[string]Script `yaml:"scripts"`

	Directory string
}

type Package struct {
	PackageManagers map[string]string       `yaml:"packageManagers"`
	InstallScopes   map[string]InstallScope `yaml:"installScopes"`
}

type InstallScope struct {
	DependsOn string              `yaml:"dependsOn"`
	Packages  map[string][]string `yaml:"packages"`
	Scripts   []string            `yaml:"scripts"`
}

type Script struct {
	Exec  string `yaml:"exec"`
	Shell string `yaml:"shell"`
}
