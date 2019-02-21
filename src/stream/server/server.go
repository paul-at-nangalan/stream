package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"os/user"
	"net/url"
	"html/template"
	"path/filepath"
	"os"
	"bytes"
	"stream/downloader/downloadpool"
	"stream/cfg"
)

type ServerCfg struct{
	cert string
	key string
	port string
}

func (p *ServerCfg)Set(){
	c := cfg.Read("server")
	p.cert = c["cert"].(string)
	p.key = c["key"].(string)
	p.port = c["port"].(string)
}

type HttpServer struct {
	svrcfg ServerCfg
	downloadpool *downloadpool.DownloadPool
}

func NewHttpServer(staging string, proxy string)*HttpServer{
	pool := downloadpool.NewDownloadPool(staging, proxy)
	return &HttpServer{
		downloadpool:pool,
	}
}

func (p *HttpServer)getTorrentFiles(dir, ext string) []string{
	files := make([]string, 0, 0)
	filelist, err := ioutil.ReadDir(dir)
	if err != nil {
		panic("Failed to get contents of " + dir)
	}
	for _, file := range filelist {
		name := file.Name()
		istorrent := strings.HasSuffix(name, ext)//".torrent")
		if istorrent {
			files = append(files, file.Name())
		}
	}
	return files
}
func (p *HttpServer)formatTorrentfiles(files []string, action, path string) string {

	html := `<table class="table">
			<thead>
				<tr>
				<th>Name</th>
				<th>Link</th>
				</tr>
			</thead>
			<tbody>`
	for _, file := range files {
		html += `<tr>
				<td>`
		html += file
		html += `	</td>`
		html += `	<td>
				<a href="`
		urlencoded := url.QueryEscape(file)
		html += "/" + path + "?file=" + urlencoded
		html += `">` + action + `</a>
				</td>
			</tr>`
	}

	html += `</tbody>
		</table>`
	return html
}
func (p *HttpServer)ListHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("List handler for " + req.URL.Path)

	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error())
	}

	templatefile, _ :=  template.ParseFiles(dir + "/views/list.html")

	params := make(map[string]interface{})
	files := p.getTorrentFiles(usr.HomeDir + "/Downloads", ".torrent")
	html := p.formatTorrentfiles(files, "Download", "download")
	params["Files"] = template.HTML(html)
	files = p.getTorrentFiles(usr.HomeDir + "/Downloads", ".inprogress")
	html = p.formatTorrentfiles(files, "--", "list")
	params["Inprogress"] = template.HTML(html)
	files = p.getTorrentFiles(usr.HomeDir + "/mov", "")
	html = p.formatTorrentfiles(files, "Play", "play")
	params["Complete"] = template.HTML(html)
	buff := bytes.NewBufferString("")
	err = templatefile.Execute(buff, params)
	if err != nil {
		log.Println("Error exporting: " + err.Error())
		panic("Failed to show view")
	}

	w.Write(buff.Bytes())
}
func (p *HttpServer)DownloadHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("Download handler for " + req.URL.Path)
	err := req.ParseForm()
	if err != nil {
		panic("Failed to parse request")
	}
	file := req.Form["file"][0]
	dl := downloadpool.Download{
		Url: file,
	}
	p.downloadpool.Enque(dl)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Download started"))
}

func (p *HttpServer)CreateServer( servermode bool) {
	p.svrcfg.Set()
	http.HandleFunc("/", p.ListHandler)
	http.HandleFunc("/list", p.ListHandler)
	http.HandleFunc("/download", p.DownloadHandler)
	if !servermode {
		http.HandleFunc("/play", PlayHandler)
		http.HandleFunc("/pause", PauseHandler)
		http.HandleFunc("/stop", QuitHandler)
		http.HandleFunc("/resume", ResumeHandler)
	}
	err := http.ListenAndServeTLS(":" + p.svrcfg.port, p.svrcfg.cert, p.svrcfg.key, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

