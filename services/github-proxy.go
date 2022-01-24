package services

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sirupsen/logrus"
)

type reverseProxyHandler struct {
	hosts   []string
	target  string
	proxies map[string]*httputil.ReverseProxy
}

func (gp *reverseProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy, ok := gp.proxies[r.Host]
	if !ok {
		targetURL := fmt.Sprintf("%s/https://%s", gp.target, r.Host)
		target, err := url.Parse(targetURL)
		if err != nil {
			logrus.Errorf("parse %s error: %v", targetURL, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(http.StatusText(http.StatusBadRequest)))
			return
		}

		proxy = httputil.NewSingleHostReverseProxy(target)
		director := proxy.Director
		proxy.Director = func(req *http.Request) {
			director(req)
			req.Host = req.URL.Host
		}
		gp.proxies[r.Host] = proxy
	}
	logrus.Infof("requesting url: %s", fmt.Sprintf("https://%s%s", r.Host, r.RequestURI))
	proxy.ServeHTTP(w, r)
}

func NewReverseProxyHandler(hosts []string, target string) *reverseProxyHandler {
	gp := &reverseProxyHandler{
		hosts:   hosts,
		target:  target,
		proxies: map[string]*httputil.ReverseProxy{},
	}
	return gp
}
