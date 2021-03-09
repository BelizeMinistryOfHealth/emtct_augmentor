package config

import (
	"testing"
)

func TestReadConf(t *testing.T) {
	_, err := ReadConf("cnf_test.yaml")
	if err != nil {
		t.Errorf("error reading configuration file: %+v", err)
	}

}
