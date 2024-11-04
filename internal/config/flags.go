package config

import "flag"

var NetAddress struct {
	StartServer string
	ResultAddr  string
}

func GetNetAddr() {
	flag.StringVar(&NetAddress.StartServer, "a", ":8080", "addr for start a server")
	flag.StringVar(&NetAddress.ResultAddr, "b", "http://localhost:8080/", "addr for base result URL")
}
