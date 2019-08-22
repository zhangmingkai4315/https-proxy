package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	log.Printf("proxy local url %s", req.RequestURI)
	proxyconfig, ok := p.config.searchProxy[req.RequestURI]
	if ok != true {
		http.Error(w, "not upstream configuration", http.StatusInternalServerError)
		return
	}
	log.Printf("find proxy setting for local url %s", req.RequestURI)
	url := fmt.Sprintf("%s%s", proxyconfig.Upstream, req.RequestURI)

	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: proxyconfig.tlsConfig,
		},
	}
	response, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.Write([]byte("proxy error:" + err.Error()))
		return
	}
	for h, val := range response.Header {
		w.Header().Set(h, strings.Join(val, " "))
	}
	w.Write(data)
}
