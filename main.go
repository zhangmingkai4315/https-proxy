package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	config string
	port   string
	help   bool
)

func init() {
	flag.StringVar(&config, "c", "config.json", "config file for application and proxy")
	flag.BoolVar(&help, "h", false, "help")
}

func main() {
	flag.Parse()
	if help == true {
		flag.Usage()
		os.Exit(0)
	}
	config, err := LoadConfig(config)
	if err != nil {
		log.Panicf("read config file error:%s", err)
	}
	proxy := NewProxy(config)
	listenAt := config.Application.ListenAt()
	log.Printf("start proxy serve in %s", listenAt)
	err = http.ListenAndServe(listenAt, proxy)
	if err != nil {
		log.Panic(err)
	}
}
