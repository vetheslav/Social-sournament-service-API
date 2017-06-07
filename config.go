package main

import (
	"github.com/zpatrick/go-config"
	"path/filepath"
	"runtime"
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
	_, b, _, _ := runtime.Caller(0)
	basePath   := filepath.Dir(b)
	dir, err := filepath.Abs(basePath)
	CheckError(err, "Open config path")
	//dir = "/Users/vetheslav/Dropbox/programming/test/Social-sournament-service-API/"

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
