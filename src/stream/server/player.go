package server

import (
	"net/http"
	"os/exec"
	"os/user"
	"io"
	"html/template"
	"bytes"
	"log"
	"path/filepath"
	"os"
)

var running *exec.Cmd = nil

func sendStdin(ch string){
	stdin, err := running.StdinPipe()
	if err != nil {
		panic("Failed to get stdin")
	}
	defer stdin.Close()
	io.WriteString(stdin, ch)
}

func displayPlayer(w http.ResponseWriter) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error())
	}
	templatefile, _ :=  template.ParseFiles(dir + "/views/list.html")
	buff := bytes.NewBufferString("")
	params := make(map[string]string)
	err = templatefile.Execute(buff, params)
	if err != nil {
		log.Println("Error exporting: " + err.Error())
		panic("Failed to show view")
	}

	w.Write(buff.Bytes())
}
func findMp4(dir string)string{
	files := getTorrentFiles(dir, ".mp4")
	if len(files) > 0 {
		return files[0]
	}
	panic("No mp4 found")
}

func PlayHandler(w http.ResponseWriter, req *http.Request) {

	if running != nil {
		sendStdin("q")
	}
	err := req.ParseForm()
	if err != nil {
		panic("Failed to parse request")
	}
	file := req.Form["file"][0]
	start := "0"
	if _, ok := req.Form["startpos"]; ok {
		start = req.Form["startpos"][0]
	}
	usr, err := user.Current()
	if err != nil {
		panic( err.Error() )
	}
	fullfile := usr.HomeDir + "/" + file
	cmd := exec.Command("omxplayer", "--pos", start,fullfile)
	err = cmd.Run()
	if err != nil {
		panic( err.Error() )
	}
	running = cmd
	sendStdin("x")
	displayPlayer(w)
}

func PauseHandler(w http.ResponseWriter, req *http.Request) {
	if running == nil {
		panic("Not running")
	}
	sendStdin("p")
	displayPlayer(w)
}
func QuitHandler(w http.ResponseWriter, req *http.Request) {
	if running == nil {
		panic("Not running")
	}
	sendStdin("q")
	displayPlayer(w)
}
func ResumeHandler(w http.ResponseWriter, req *http.Request) {
	if running == nil {
		panic("Not running")
	}
	sendStdin(" ")
	displayPlayer(w)
}
