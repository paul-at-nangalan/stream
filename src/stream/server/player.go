package server

import (
	"net/http"
	"os/exec"
	"os/user"
)

func PlayHandler(w http.ResponseWriter, req *http.Request) {

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
}
