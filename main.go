package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Configuration struct {
	Playlists []PlaylistConf
}

type PlaylistConf struct {
	Type string
	Name string
	IP   string
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

func handler(w http.ResponseWriter, r *http.Request) {
	// read configuration
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, plConf := range configuration.Playlists {
		if plConf.Type == "pomoyka" {
			if r.URL.Path[1:] == plConf.Name {
				resp, err := http.Get("http://pomoyka.win/trash/ttv-list/ttv.all.iproxy.xspf?ip=" + plConf.IP + ":6878")
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

				//xmlFile, err := os.Open("ttv.all.iproxy.xspf")
				//if err != nil {
				//fmt.Println(err)
				//}

				//// defer the closing of our xmlFile so that we can parse it later on
				//defer xmlFile.Close()
				//xmlByte, err := ioutil.ReadAll(xmlFile)
				//if err != nil {
				//fmt.Errorf("Read xml: %v", err)
				//}

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
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
