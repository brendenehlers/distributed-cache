package main

import (
	"flag"
	"fmt"

	regmap "github.com/brendenehlers/go-distributed-cache/registry-node/map"
	"github.com/brendenehlers/go-distributed-cache/registry-node/server"
)

var (
	hostnameFlag = flag.String("hostname", "localhost", "host name for the server")
	portFlag     = flag.Int("port", 8081, "port for the server")
)

func init() {
	flag.Parse()
}

func main() {
	host := fmt.Sprintf("%s:%d", *hostnameFlag, *portFlag)
	reg := regmap.New()
	server := server.New(host, reg)

	server.Start()
}
