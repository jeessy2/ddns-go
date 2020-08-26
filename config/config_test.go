package config

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	conf := &Config{}
	byt, err := ioutil.ReadFile("../config.yaml")
	if err != nil {
		t.Error(err)
	}

	yaml.Unmarshal(byt, conf)
	if "alidns" != conf.DNS.Name {
		t.Error("DNS Name必须为alidns")
	}
}
