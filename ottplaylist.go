package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Configuration struct {
	Playlists []PlaylistConf
	Port      int
}

type PlaylistConf struct {
	Type     string
	Playlist string
	Name     string
	IP       string
	Port     int
}

type Playlist struct {
	XMLName xml.Name `xml:"playlist"`
	Title   string   `xml:"title"`
	Groups  []Group  `xml:"extension>node"`
	Tracks  []Track  `xml:"trackList>track"`
}

type Group struct {
	XMLName  xml.Name  `xml:"node"`
	Title    string    `xml:"title,attr"`
	Channels []Channel `xml:"item"`
}

type Channel struct {
	XMLName xml.Name `xml:"item"`
	Id      string   `xml:"tid,attr"`
}

type Track struct {
	XMLName  xml.Name `xml:"track"`
	Id       string   `xml:"extension>id"`
	Location string   `xml:"location"`
	Title    string   `xml:"title"`
}

func (conf *Configuration) handler(w http.ResponseWriter, r *http.Request) {
	for _, plConf := range conf.Playlists {
		if plConf.Type == "pomoyka" {
			if r.URL.Path[1:] == plConf.Name {
				url := "http://pomoyka.win/trash/ttv-list/"
				if plConf.Playlist == "proxy" {
					url += "ttv.all.proxy.xspf"
				} else {
					url += "ttv.all.iproxy.xspf"
				}
				resp, err := http.Get(url + "?ip=" + plConf.IP + ":" + strconv.Itoa(plConf.Port))
				if err != nil {
					fmt.Errorf("GET error: %v", err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Errorf("Status error: %v", resp.StatusCode)
				}

				xmlByte, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Errorf("Read body: %v", err)
				}
				var pl Playlist
				xml.Unmarshal(xmlByte, &pl)
				var ch = make(map[string]map[string]string)
				for i := 0; i < len(pl.Tracks); i++ {
					var data = make(map[string]string)
					data["title"] = pl.Tracks[i].Title
					data["location"] = pl.Tracks[i].Location
					ch[pl.Tracks[i].Id] = data
				}
				fmt.Fprintf(w, "#EXTM3U\n")
				for i := 0; i < len(pl.Groups); i++ {
					for j := 0; j < len(pl.Groups[i].Channels); j++ {
						fmt.Fprintf(w, "#EXTINF:0,"+ch[pl.Groups[i].Channels[j].Id]["title"]+"\n")
						fmt.Fprintf(w, "#EXTGRP:"+pl.Groups[i].Title+"\n")
						fmt.Fprintf(w, ch[pl.Groups[i].Channels[j].Id]["location"]+"\n")
					}
				}
			}
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
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(configuration.Port), nil))
}
