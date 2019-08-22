package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type appConfig struct {
	Server string `json:"server"`
	Port   string `json:"port"`
}

func (app *appConfig) ListenAt() string {
	return app.Server + ":" + app.Port
}

type proxyConfig struct {
	Location   string `json:"location"`
	Upstream   string `json:"upstream"`
	ClientCert string `json:"client_cert"`
	ClientKey  string `json:"client_key"`
	CACert     string `json:"ca_cert"`
	tlsConfig  *tls.Config
}

// Config hold all configuration
type Config struct {
	Application appConfig      `json:"app"`
	ProxyList   []*proxyConfig `json:"proxy"`
	searchProxy map[string]*proxyConfig
}

// LoadConfig read config file and
// return a config struct
func LoadConfig(configfile string) (*Config, error) {
	jsonFile, err := os.Open(configfile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteJSONData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(byteJSONData, config)
	if err != nil {
		return nil, err
	}
	searchProxy := make(map[string]*proxyConfig)
	for _, proxy := range config.ProxyList {
		err = loadCertAndKeyFile(proxy)
		if err != nil {
			return nil, err
		}
		searchProxy[proxy.Location] = proxy

	}
	config.searchProxy = searchProxy
	return config, nil
}

func loadCertAndKeyFile(proxy *proxyConfig) error {
	var cert tls.Certificate
	var err error
	if proxy.ClientCert != "" && proxy.ClientKey != "" {
		log.Printf("load client cert for proxy:%s", proxy.Location)
		cert, err = tls.LoadX509KeyPair(proxy.ClientCert, proxy.ClientKey)
		if err != nil {
			return err
		}
	}
	clientCACertPool := x509.NewCertPool()
	if proxy.CACert != "" {
		log.Printf("load ca cert for proxy:%s", proxy.Location)
		caCert, err := ioutil.ReadFile(proxy.CACert)
		if err != nil {
			return err
		}
		ok := clientCACertPool.AppendCertsFromPEM(caCert)
		if ok != true {
			return errors.New("ca certs add fail")
		}
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      clientCACertPool,
	}
	proxy.tlsConfig = tlsConfig
	return nil
}
