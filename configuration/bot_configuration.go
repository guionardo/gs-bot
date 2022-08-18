package configuration

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Configuration struct {
	Repository RepositoryConfiguration
}

var configuration *Configuration

func GetConfiguration() (cfg Configuration, err error) {
	if configuration == nil {
		ctx := context.Background()
		configuration = &Configuration{}
		err = envconfig.Process(ctx, configuration)
		configuration.FixDefaults()
	}
	return *configuration, err
}

func (c *Configuration) FixDefaults() {
	c.Repository.FixDefaults()
}
