package cfg

import (
	"io/ioutil"
	"path"
	"log"
	"encoding/json"
)

var cfgpath string

func Set(path string){
	cfgpath = path
}

func Read(what string)map[string]interface{}{
	m := make(map[string]interface{})
	file := path.Join(cfgpath, what + ".json")
	cfgdata, err := ioutil.ReadFile(file)
	if err != nil{
		log.Panic("Failed to read cfg ", what)
	}
	err = json.Unmarshal(cfgdata, m)
	if err != nil {
		log.Panic("failed to read cfg ", err)
	}
	return m
}
