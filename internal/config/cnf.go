package config

import "github.com/spf13/viper"
import log "github.com/sirupsen/logrus"

type DbConf struct {
	Username string
	Password string
	Database string
	Host     string
}

type AuthConf struct {
	JwkUrl   string
	Issuer   string
	Audience string
}

type AppConf struct {
	EmtctDb   DbConf
	Auth      AuthConf
	AcsisDb   DbConf
	ProjectId string
}

type Firebase struct {
	ProjectId string
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

	var c DbConf
	err := viper.Sub("emtct_db").Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	var a AuthConf
	err = viper.Sub("auth").Unmarshal(&a)
	if err != nil {
		return nil, err
	}

	var acsisConf DbConf
	err = viper.Sub("acsis_db").Unmarshal(&acsisConf)
	if err != nil {
		return nil, err
	}

	var firebase Firebase
	err = viper.Sub("firebase").Unmarshal(&firebase)
	if err != nil {
		return nil, err
	}

	appConf := AppConf{
		EmtctDb:   c,
		Auth:      a,
		AcsisDb:   acsisConf,
		ProjectId: firebase.ProjectId,
	}

	return &appConf, nil
}
