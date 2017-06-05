package main

import (
	"github.com/zpatrick/go-config"
	"os"
	"path/filepath"
)

var conf Config

type Config struct {
	configFile *config.Config
	mysql      ConfigMysql
	httpServer HttpServer

}

func (conf *Config) initConfig() {
	conf.getConfigObject()
	conf.parseConfigVariables()
}

func (conf *Config) getConfigObject() {
	configPath := conf.openConfigPath()
	conf.openConfigFile(configPath)
}

func (conf *Config) openConfigPath() (dir string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	CheckError(err, "Open config path")
	//TODO убрать абсолютный путь
	dir = "/Users/vetheslav/Dropbox/programming/test/Social-sournament-service-API/"

	return
}

func (conf *Config) openConfigFile(dir string) {
	yamlFile := config.NewYAMLFile(dir + "/config.yaml")
	conf.configFile = config.NewConfig([]config.Provider{yamlFile})
	err := conf.configFile.Load()
	CheckError(err, "Load config file")
}

func (conf *Config) parseConfigVariables() {
	conf.parseConfigMysql()
	conf.parseConfigHttp()
}

func (conf *Config) parseConfigMysql() {
	conf.mysql = ConfigMysql{}
	conf.mysql.parseConfigMysqlHost()
	conf.mysql.parseConfigMysqlPort()
	conf.mysql.parseConfigMysqlUser()
	conf.mysql.parseConfigMysqlPassword()
	conf.mysql.parseConfigMysqlTable()
}

func (conf *Config) parseConfigHttp() {
	conf.httpServer = HttpServer{}
	conf.httpServer.parseConfigHttpHost()
	conf.httpServer.parseConfigHttpPort()
}