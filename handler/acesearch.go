package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"ottplaylist/types"
	"strconv"
	"time"
)

// acesearch channels final structure
type AceChan map[string]AceGetChannel

// acesearch channels previous cache structure
type AcePrev map[string]AceGetChannel

type AceGetChannel struct {
	Name     string
	Avail    float64
	Upd      int64
	Cat      string
	Infohash string
	T        int64
}

// channels from acesearch
type Ace []struct {
	Infohash                string
	Name                    string
	Availability            float64
	Availability_updated_at int64
	Categories              []string
}

const AVAIL_THRESHOLD = (8 * 86400)
const CHANNEL_EXPIRED = 86400

func GetAcePlaylist(link string, IP string, port uint16) types.Channels {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Errorf("GET json error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	ace := Ace{}
	err = decoder.Decode(&ace)
	if err != nil {
		fmt.Println("error: Bad Json, ", err)
	}
	if len(ace) < 100 {
		panic("Too few channels")
	}
	prevPath := "tmp/"
	prevFileLink := prevPath + "ace.prev.json"
	prevFile, err := os.Open(prevFileLink)
	if err != nil && prevFile != nil {
		panic(err)
	}
	if prevFile != nil {
		defer func() {
			if err := prevFile.Close(); err != nil {
				panic(err)
			}
		}()
	}
	acePrev := make(AcePrev)
	if prevFile != nil {
		decoder = json.NewDecoder(prevFile)
		err = decoder.Decode(&acePrev)
		if err != nil {
			fmt.Println("error: Bad prev Json, ", err)
		}
	}

	timeNow := time.Now().Unix()

	aceChannels := make(AceChan)
	for _, aceCh := range ace {
		updt := timeNow - aceCh.Availability_updated_at
		if aceCh.Availability < 0.8 || updt > AVAIL_THRESHOLD {
			continue
		}
		val, ok := aceChannels[aceCh.Name]
		if !ok || val.Upd < aceCh.Availability_updated_at {
			cat := "none"
			if len(aceCh.Categories) > 0 {
				cat = aceCh.Categories[0]
			}
			aceChannels[aceCh.Name] = AceGetChannel{
				Name:     aceCh.Name,
				Avail:    aceCh.Availability,
				Upd:      aceCh.Availability_updated_at,
				Cat:      cat,
				Infohash: aceCh.Infohash,
				T:        timeNow,
			}
			fmt.Printf("adding search channel \"%s\" (infohash %s upd %d)\n", aceCh.Name, aceCh.Infohash, updt)
		}
	}
	prevF, _ := json.MarshalIndent(aceChannels, "", "  ")
	err = ioutil.WriteFile(prevFileLink, prevF, 0644)
	if err != nil {
		fmt.Errorf("Error while write cache file: %v\n", err)
	}
	// add channels from previous search cache
	for _, p := range acePrev {
		age := timeNow - p.T
		updt := timeNow - p.Upd
		if age < CHANNEL_EXPIRED && updt < AVAIL_THRESHOLD {
			_, ok := aceChannels[p.Name]
			if !ok {
				aceChannels[p.Name] = p
				fmt.Printf("adding previous search channel \"%s\" (infohash %s age %d upd %d)\n", p.Name, p.Infohash, age, updt)
			} else if aceChannels[p.Name].Upd < p.Upd {
				aceChannels[p.Name] = p
				fmt.Printf("replacing search channel \"%s\"\n", p.Name)
				if aceChannels[p.Name].Infohash != p.Infohash {
					fmt.Printf("infohash mismatch for channel \"%s\" (infohash %s => %s age %d upd %d)\n", p.Name, aceChannels[p.Name].Infohash, p.Infohash, age, updt)
				}
			}
		}
	}
	channels := types.Channels{}
	for _, c := range aceChannels {
		channels = append(channels, types.Channel{
			Name:     c.Name,
			URL:      "http://" + IP + ":" + strconv.Itoa(int(port)) + "/ace/getstream?infohash=" + c.Infohash,
			Category: c.Cat,
		})
	}
	return channels
}
