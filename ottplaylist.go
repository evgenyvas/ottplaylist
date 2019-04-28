package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"ottplaylist/format"
	"ottplaylist/handler"
	"ottplaylist/types"
	"strconv"
)

type Configuration struct {
	Types map[string]struct {
		Link    string
		Handler string
	}
	Playlists []struct {
		Type   string
		Name   string
		IP     string
		Port   uint16
		Format string
	}
	Port uint16
}

func (conf *Configuration) handler(w http.ResponseWriter, r *http.Request) {
	for _, plConf := range conf.Playlists {
		if r.URL.Path[1:] == plConf.Name { // check route
			var ch types.Channels
			if conf.Types[plConf.Type].Handler == "pomoyka.xspf" {
				ch = handler.GetPomoykaPlaylist(conf.Types[plConf.Type].Link, plConf.IP, plConf.Port)
			} else if conf.Types[plConf.Type].Handler == "acesearch" {
				ch = handler.GetAcePlaylist(conf.Types[plConf.Type].Link, plConf.IP, plConf.Port)
			}
			var pl string
			if plConf.Format == "ott.m3u" {
				pl = format.M3U(ch)
			}
			//fmt.Printf("%+v\n", ch)
			fmt.Fprintf(w, pl)
		}
	}
}

func main() {
	// read configuration
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	http.HandleFunc("/", configuration.handler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(int(configuration.Port)), nil))
}
