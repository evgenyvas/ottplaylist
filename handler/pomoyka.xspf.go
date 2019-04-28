package handler

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"ottplaylist/types"
	"strconv"
)

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

func GetPomoykaPlaylist(link string, IP string, port uint16) types.Channels {
	channels := types.Channels{}
	resp, err := http.Get(link + "?ip=" + IP + ":" + strconv.Itoa(int(port)))
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
	for i := 0; i < len(pl.Groups); i++ {
		for j := 0; j < len(pl.Groups[i].Channels); j++ {
			channels = append(channels, types.Channel{
				Name:     ch[pl.Groups[i].Channels[j].Id]["title"],
				URL:      ch[pl.Groups[i].Channels[j].Id]["location"],
				Category: pl.Groups[i].Title,
			})
		}
	}
	return channels
}
