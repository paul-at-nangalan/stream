package server

import (
	"io/ioutil"
	"fmt"
	"log"
	"net/http"
	"strings"
	"os/user"
	"net/url"
	"html/template"
	"cmd/vet/internal/cfg"
	"path/filepath"
	"os"
	"bytes"
)

func getTorrentFiles(dir string) []string{
	files := make([]string, 0, 0)
	filelist, err := ioutil.ReadDir(dir)
	if err != nil {
		panic("Failed to get contents of " + dir)
	}
	for _, file := range filelist {
		name := file.Name()
		istorrent := strings.HasSuffix(name, ".torrent")
		if istorrent {
			files = append(files, file)
		}
	}
	return files
}
func askTorrent() (torrentfile string) {
	fmt.Println("Please select: ")
	files := getTorrentFiles("/home/pi/Downloads")
	indx := 1
	for _, file := range files {
		fmt.Println(indx, ":", file)
	}

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

	files := getTorrentFiles(usr.HomeDir() + "/Downloads")
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
		urlencoded := url.PathEscape(file)
		html += urlencoded
		html += `">Start download</a>
				</td>
			</tr>`
	}

	html += `</tbody>
		</table>`
	buff := bytes.NewBufferString("")
	params := make(map[string]string)
	params["Files"] = html
	err = templatefile.Execute(buff, params)
	if err != nil {
		log.Println("Error exporting: " + err.Error())
		panic("Failed to export to pdf")
	}

	w.Write(buff.Bytes())
}
func DownloadHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("Download handler for " + req.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Page not found"))
}
func GeneralHandler(w http.ResponseWriter, req *http.Request) {
	log.Print("General handler for " + req.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Page not found"))
}

func CreateServer() {
	http.HandleFunc("/", GeneralHandler)
	http.HandleFunc("/list", ListHandler)
	http.HandleFunc("/download", DownloadHandler)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

