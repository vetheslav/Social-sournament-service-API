package main

import (
	"os"
	"reflect"
	"testing"
)

func openConfigFile() {
	conf.openConfigFile(os.Getenv("GOPATH"))
}

func TestOpenConfigFile(t *testing.T) {
	openConfigFile()

	confObjectType := reflect.TypeOf(conf.configFile).String()
	checkStrings(t, confObjectType, "*config.Config")
}

func TestParseConfigMysql(t *testing.T) {
	openConfigFile()

	conf.parseConfigMysql()
	checkStrings(t, conf.mysql.host, "127.0.0.1")
	checkStrings(t, conf.mysql.user, "vista")
	checkStrings(t, conf.mysql.password, "vista123")
	checkStrings(t, conf.mysql.table, "social_tourment")
	checkStrings(t, conf.mysql.port, "3306")
}

func TestParseConfig(t *testing.T) {
	openConfigFile()

	conf.parseConfigVariables()
	checkStrings(t, conf.mysql.host, "127.0.0.1")
	checkStrings(t, conf.httpServer.port, "1323")
}

func checkStrings(t *testing.T, expected string, got string) {
	if expected != got {
		t.Error(
			"Expected", expected,
			"got", got,
		)
	}
}
