package main

import (
	"log"
	"runtime/debug"
	"os"
	"stream/downloader"
	"stream/server"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}

	staging := usr.HomeDir + "/mov"
	defer func(){
		if r := recover(); r != nil {
			log.Println(r)
			debug.PrintStack()
		}
	}()
	os.MkdirAll(staging, os.ModeDir | os.ModePerm)

	dl := downloader.New(staging)
	server.CreateServer(&dl, os.Args[1])
}
