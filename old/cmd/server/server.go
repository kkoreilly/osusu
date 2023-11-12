package main

import (
	"github.com/kkoreilly/osusu/client"
	"github.com/kkoreilly/osusu/server"
)

func main() {
	// need to also start client on server so crawlers can crawl for SEO
	client.Start()
	server.Start()
}
