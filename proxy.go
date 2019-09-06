package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Proxy define the proxy behaivor struct
type Proxy struct {
	config *Config
}

// NewProxy create a new proxy handler
func NewProxy(config *Config) *Proxy {
	return &Proxy{
		config: config,
	}
}

// func (p *Proxy) directProxy()
func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	log.Infof("proxy local url %s", req.RequestURI)
	proxyconfig, ok := p.config.searchProxy[req.RequestURI]
	if ok != true {
		log.Errorf("no upstream configuration for %s", req.RequestURI)
		http.Error(w, "not upstream configuration", http.StatusInternalServerError)
		return
	}

	client, ok := p.config.client[req.RequestURI]
	if ok != true {
		log.Errorf("no upstream client for %s", req.RequestURI)
		http.Error(w, "not upstream client available", http.StatusInternalServerError)
		return
	}

	log.Debugf("==>receive request from %s", req.RequestURI)
	proxyReq, err := http.NewRequest(req.Method, proxyconfig.Upstream, bytes.NewReader(body))
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}

	log.Debugf("--> %s proxy to %s", req.RequestURI, proxyconfig.Upstream)
	response, err := client.Do(proxyReq)
	if err != nil {
		log.Errorf("--> request to %s error: %s took %s", proxyconfig.Upstream, err.Error(), time.Since(start))
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	log.Debugf("<-- receive data %s success took %s", proxyconfig.Upstream, time.Since(start))
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Errorf("<-- read response from %s error: %s", proxyconfig.Upstream, err.Error())
		w.Write([]byte("proxy error:" + err.Error()))
		return
	}
	for h, val := range response.Header {
		w.Header().Set(h, strings.Join(val, " "))
	}
	log.Debugf("<== send data to %s took %s", req.RequestURI, time.Since(start))
	w.Write(data)

}
