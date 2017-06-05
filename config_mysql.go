package main

type ConfigMysql struct {
	host     string
	port     string
	user     string
	password string
	table    string
}

func (mysql *ConfigMysql) parseConfigMysqlHost() {
	var err error
	mysql.host, err = conf.configFile.String("mysql.host")
	CheckError(err, "Not found mysql.host")
}

func (mysql *ConfigMysql) parseConfigMysqlPort() {
	var err error
	mysql.port, err = conf.configFile.String("mysql.port")
	if err != nil {
		mysql.port = "3306"
	}
}

func (mysql *ConfigMysql) parseConfigMysqlUser() {
	var err error
	mysql.user, err = conf.configFile.String("mysql.user")
	CheckError(err, "Not found mysql.user")
}

func (mysql *ConfigMysql) parseConfigMysqlPassword() {
	var err error
	mysql.password, err = conf.configFile.String("mysql.pass")
	CheckError(err, "Not found mysql.pass")
}

func (mysql *ConfigMysql) parseConfigMysqlTable() {
	var err error
	mysql.table, err = conf.configFile.String("mysql.table")
	CheckError(err, "Not found mysql.table")
}
