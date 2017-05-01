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
	"stream/downloader"
)

var dlhandler *downloader.Downloader

func getTorrentFiles(dir, ext string) []string{
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
func formatTorrentfiles(files []string, action, path string) string {

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
func ListHandler(w http.ResponseWriter, req *http.Request) {
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
	files := getTorrentFiles(usr.HomeDir + "/Downloads", ".torrent")
	html := formatTorrentfiles(files, "Download", "download")
	params["Files"] = template.HTML(html)
	files = getTorrentFiles(usr.HomeDir + "/Downloads", ".inprogress")
	html = formatTorrentfiles(files, "--", "list")
	params["Inprogress"] = template.HTML(html)
	files = getTorrentFiles(usr.HomeDir + "/mov", "")
	html = formatTorrentfiles(files, "Play", "play")
	params["Complete"] = template.HTML(html)
	buff := bytes.NewBufferString("")
	err = templatefile.Execute(buff, params)
	if err != nil {
		log.Println("Error exporting: " + err.Error())
		panic("Failed to show view")
	}

	w.Write(buff.Bytes())
}
func DownloadHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("Download handler for " + req.URL.Path)
	err := req.ParseForm()
	if err != nil {
		panic("Failed to parse request")
	}
	file := req.Form["file"][0]
	go dlhandler.Start(file)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Download started"))
}

func CreateServer(dl *downloader.Downloader, port string) {
	dlhandler = dl
	http.HandleFunc("/", ListHandler)
	http.HandleFunc("/list", ListHandler)
	http.HandleFunc("/download", DownloadHandler)
	http.HandleFunc("/play", PlayHandler)
	http.HandleFunc("/pause", PauseHandler)
	http.HandleFunc("/stop", QuitHandler)
	http.HandleFunc("/resume", ResumeHandler)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

