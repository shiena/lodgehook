package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/shiena/lodgehook"
	"github.com/shiena/lodgehook/hook"
)

func main() {
	var sourceAddress string
	var forwardURL string
	var idobataToken string
	flag.StringVar(&sourceAddress, "src", ":4000", "source address")
	flag.StringVar(&forwardURL, "forward", "http://localhost:3000", "forward URL")
	flag.StringVar(&idobataToken, "idobata", "", "idobata hook api token")
	flag.Parse()

	url, err := url.Parse(forwardURL)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = lodgehook.NewLodgeHookTransport(
		hook.NewIdobataHook(idobataToken),
	)
	server := http.Server{
		Addr:    sourceAddress,
		Handler: proxy,
	}
	log.Printf("start lodge hook %s -> %s\n", sourceAddress, forwardURL)
	log.Fatal(server.ListenAndServe())
}
