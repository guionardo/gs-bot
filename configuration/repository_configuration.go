package configuration

type RepositoryConfiguration struct {
	ConnectionString string `env:"REPOSITORY_CONNECTION_STRING"`
}

func (rc *RepositoryConfiguration) FixDefaults() {
	if len(rc.ConnectionString) == 0 {
		panic("REPOSITORY_CONNECTION_STRING is required")
	}
}
