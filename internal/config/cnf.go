package config

import "github.com/spf13/viper"
import log "github.com/sirupsen/logrus"

type DbConf struct {
	DbUsername string
	DbPassword string
	DbDatabase string
	DbHost     string
}

// ReadConf reads a yaml file and unmarshalls its content.
// The yaml file must have root siblings for the following environments:
// prod, test, dev
func ReadConf(fileName, stage string) (*DbConf, error) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(fileName)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	log.Infof("using configuration file: %s", fileName)

	v := viper.Sub(stage)
	var c *DbConf
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
