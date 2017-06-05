package main

type HttpServer struct {
	host     string
	port     string
}

func (http *HttpServer) parseConfigHttpHost() {
	var err error
	http.host, err = conf.configFile.String("http_server.host")
	CheckError(err, "Not found http_server.host")
}

func (http *HttpServer) parseConfigHttpPort() {
	var err error
	http.port, err = conf.configFile.String("http_server.port")
	if err != nil {
		http.port = "80"
	}
}