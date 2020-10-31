package config

import (
	"testing"
)

func TestReadConf(t *testing.T) {
	conf, err := ReadConf("cnf_test.yaml")
	if err != nil {
		t.Errorf("error reading configuration file: %+v", err)
	}
	if conf.EmtctDb.Username != "postgres" {
		t.Errorf("want: %s got: %s", "postgres", conf.EmtctDb.Username)
	}
}
