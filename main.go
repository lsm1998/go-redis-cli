package main

import (
	"flag"
	"os"
)

var host = flag.String("h", "localhost", "redis host")
var port = flag.Int("p", 6379, "redis port")
var pass = flag.String("pass", "", "redis pass")

func main() {
	flag.Parse()
	var client = &RedisClient{
		host: *host,
		port: *port,
		pass: *pass,
	}
	if err := client.Connect(); err != nil {
		panic(err)
	}
	Start(os.Stdin, os.Stdout, client)
}
