package main

import (
	"flag"
	"log"
	"runtime/debug"
	"github.com/anacrolix/torrent"
	"os"
)

func main() {
	staging := "/home/pi/mov"
	defer func(){
		if r := recover(); r != nil {
			log.Println(r)
			debug.PrintStack()
		}
	}()
	torrentfile := ""
	flag.StringVar(&torrentfile, "torrent-file", "", "Torrent file")	
	flag.Parse()	
	os.MkdirAll(staging, os.ModeDir | os.ModePerm)
	cfg := torrent.Config{DataDir: staging}
	c, err := torrent.NewClient(&cfg)
	if err != nil {
		panic("failed to create torrent with err " + err.Error())
	}
	defer c.Close()
	t, err := c.AddTorrentFromFile(torrentfile)
	if err != nil {
		panic("failed to add torrent file to torrent with err " + err.Error())
	}
	<-t.GotInfo()
	t.DownloadAll()
	c.WaitAll()
	log.Print(torrentfile + ", torrent downloaded")

}
