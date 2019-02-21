package downloadpool

import "stream/downloader"

type Download struct{
	Url string
}

///we can't pool these because they have to bind to a port
type DownloadPool struct{
	downloader downloader.Downloader
	next chan Download
}

func NewDownloadPool(staging string, proxyurl string)*DownloadPool{
	c := make(chan Download, 1000)
	dl := downloader.New(staging, proxyurl)
	dlp := DownloadPool{
		downloader:dl,
		next:c,
	}
	go dlp.Run()

	return &dlp
}

func (p *DownloadPool)Enque(dl Download){
	p.next <- dl
}


func (p *DownloadPool)Run(){
	for download := range p.next{
		p.downloader.Start(download.Url)
	}
}
