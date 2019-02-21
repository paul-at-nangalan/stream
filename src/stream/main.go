package main

import (
	"log"
	"runtime/debug"
	"os"
	"stream/downloader"
	"stream/server"
	"os/user"
	"stream/downloader/downloadpool"
	"flag"
	"stream/cfg"
)

func main() {

	defer func(){
		if r := recover(); r != nil {
			log.Println(r)
			debug.PrintStack()
		}
	}()

	cfgdir := "../cfg"
	flag.StringVar(&cfgdir, "cfg", "../cfg", "Config dir")
	flag.Parse()

	cfg.Set(cfgdir)
	cfgmap := cfg.Read("general")
	staging := os.ExpandEnv(cfgmap["staging"].(string))
	os.MkdirAll(staging, os.ModeDir | os.ModePerm)
	proxyurl := os.ExpandEnv(cfgmap["proxyurl"].(string))

	svr := server.NewHttpServer(staging, proxyurl)
	isserver := false
	if svr, ok := cfgmap["isserver"]; ok {
		isserver = svr.(bool)
	}
	svr.CreateServer(isserver)
}
