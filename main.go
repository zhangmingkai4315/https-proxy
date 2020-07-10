package main

import (
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router/gin"
	"flag"
	"log"
	"os"
)

func main(){
	port := flag.Int("p", 0, "port of the service")
	logLevel := flag.String("l", "ERROR", "logging level")
	debug := flag.Bool("d", false, "Enable the debug")
	configFile := flag.String("c", "krakend.json","path to the configuration file")
	flag.Parse()

	parser := config.NewParser()
	serviceConfig, err := parser.Parse(*configFile)
	if err != nil{
		log.Fatal("Error:", err.Error())
	}
	serviceConfig.Debug = serviceConfig.Debug || *debug
	if *port != 0{
		serviceConfig.Port = *port
	}
	logger, _ := logging.NewLogger(*logLevel, os.Stdout, "[Proxy]")
	routerFactory := gin.DefaultFactory(proxy.DefaultFactory(logger), logger)
	routerFactory.New().Run(serviceConfig)
}