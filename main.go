package main

import (
	"flag"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	config string
	port   string
	help   bool
	debug  bool
)

func init() {
	flag.BoolVar(&debug, "d", false, "enable debug mode")
	flag.StringVar(&config, "c", "config.json", "config file for application and proxy")
	flag.BoolVar(&help, "h", false, "help")
}

func main() {
	flag.Parse()
	if help == true {
		flag.Usage()
		os.Exit(0)
	}
	if debug == true {
		log.Info("set application in debug mode")
		log.SetLevel(log.DebugLevel)
	}
	config, err := LoadConfig(config)
	if err != nil {
		log.Panicf("read config file error:%s", err)
	}
	proxy := NewProxy(config)
	listenAt := config.Application.ListenAt()
	log.Infof("start proxy serve in %s", listenAt)
	err = http.ListenAndServe(listenAt, proxy)
	if err != nil {
		log.Error(err)
	}
}
