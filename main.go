package main

import (
	"flag"
	"log"

	"github.com/panyam/blog/web"
)

var (
	gw_addr = flag.String("gw_addr", web.DefaultGatewayAddress(), "Address where the http grpc gateway endpoint is running")
)

func main() {
	flag.Parse()

	ohweb, err := web.NewWebApp(*gw_addr)
	if err != nil {
		log.Fatal(err)
	}
	ohweb.Start()
}
