package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type appConfig struct {
	Server string `json:"server"`
	Port   string `json:"port"`
}

func (app *appConfig) ListenAt() string {
	return app.Server + ":" + app.Port
}

type proxyConfig struct {
	SkipServerValidation bool   `json:"skip_server_ssL_validation"`
	Location             string `json:"location"`
	Upstream             string `json:"upstream"`
	ClientCert           string `json:"client_cert"`
	ClientKey            string `json:"client_key"`
	CACert               string `json:"ca_cert"`
	tlsConfig            *tls.Config
}

// Config hold all configuration
type Config struct {
	Application appConfig      `json:"app"`
	ProxyList   []*proxyConfig `json:"proxy"`
	searchProxy map[string]*proxyConfig
	client      map[string]*http.Client
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
	client := make(map[string]*http.Client)
	for _, proxy := range config.ProxyList {
		err = loadCertAndKeyFile(proxy)
		if err != nil {
			return nil, err
		}
		searchProxy[proxy.Location] = proxy
		client[proxy.Location] = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: proxy.tlsConfig,
			},
		}

	}
	config.searchProxy = searchProxy
	config.client = client
	return config, nil
}

func loadCertAndKeyFile(proxy *proxyConfig) error {
	var cert tls.Certificate
	var err error
	if proxy.ClientCert != "" && proxy.ClientKey != "" {
		log.Infof("load client cert for proxy:%s", proxy.Location)
		cert, err = tls.LoadX509KeyPair(proxy.ClientCert, proxy.ClientKey)
		if err != nil {
			return err
		}
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	if proxy.CACert != "" {
		log.Infof("load ca cert for proxy:%s", proxy.Location)
		caCert, err := ioutil.ReadFile(proxy.CACert)
		if err != nil {
			return err
		}
		ok := rootCAs.AppendCertsFromPEM(caCert)
		if ok != true {
			return errors.New("ca certs add fail")
		}
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            rootCAs,
		InsecureSkipVerify: proxy.SkipServerValidation,
	}
	proxy.tlsConfig = tlsConfig
	return nil
}
