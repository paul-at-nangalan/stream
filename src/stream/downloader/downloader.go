package downloader

import (
	"github.com/anacrolix/torrent"
	"os"
	"os/user"
	"log"
)

type Downloader struct {
	staging string
}

func New(staging string) Downloader {
	return Downloader{staging: staging}
}

func (p *Downloader)Start(origtorrentfile string){
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	dir := usr.HomeDir + "/Downloads"
	torrentfile := dir + "/" + origtorrentfile
	cfg := torrent.Config{DataDir: p.staging}
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
	torrentfile = dir + "/" + origtorrentfile + ".inprogress"
	origtorrentfile = dir + "/" + origtorrentfile
	os.Rename(origtorrentfile, torrentfile)
	c.WaitAll()
	os.Rename(torrentfile, origtorrentfile + ".complete")
}
