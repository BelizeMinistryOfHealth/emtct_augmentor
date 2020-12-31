package config

import "github.com/spf13/viper"
import log "github.com/sirupsen/logrus"

type AppConf struct {
	ProjectId       string
	Port            int
	FirestoreApiKey string
}

type Firebase struct {
	ProjectId string
	ApiKey    string
}

// ReadConf reads a yaml file and unmarshalls its content.
// The yaml file must have root siblings for the following environments:
// prod, test, dev
func ReadConf(fileName string) (*AppConf, error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(fileName)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	log.Infof("using configuration file: %s", fileName)

	var firebase Firebase
	err := viper.Sub("firebase").Unmarshal(&firebase)
	if err != nil {
		return nil, err
	}

	appConf := AppConf{
		ProjectId:       firebase.ProjectId,
		FirestoreApiKey: firebase.ApiKey,
	}

	return &appConf, nil
}
