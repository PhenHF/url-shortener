package config

import (
	"flag"
	"os"
)

func init() {
	GetNetAddr()
}

var NetAddress struct {
	StartServer string
	ResultAddr  string
}

func GetNetAddr() {
	flag.StringVar(&NetAddress.StartServer, "a", ":8080", "addr for start a server")
	flag.StringVar(&NetAddress.ResultAddr, "b", "http://localhost:8080/", "addr for base result URL")

	flag.Parse()

	if serverAddr := os.Getenv("SERVER_ADDRESS"); serverAddr != "" {
		NetAddress.StartServer = serverAddr
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		NetAddress.ResultAddr = baseURL
	}

}
