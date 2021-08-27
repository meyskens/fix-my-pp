package main

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/meyskens/recent-beater/pkg/bplist"
)

type score struct {
	Ranked bool `json:"ranked"`
	UID    int  `json:"uid"`
	Scores []struct {
		PlayerID int64         `json:"playerId"`
		Name     string        `json:"name"`
		Rank     int           `json:"rank"`
		Score    int           `json:"score"`
		Pp       float64       `json:"pp"`
		Mods     []interface{} `json:"mods"`
	} `json:"scores"`
	Mods  bool `json:"mods"`
	Valid bool `json:"valid"`
}

type data struct {
	Diffs []struct {
		Pp     string `json:"pp"`
		Star   string `json:"star"`
		Scores string `json:"scores"`
		Diff   string `json:"diff"`
		Type   int    `json:"type"`
		Len    int    `json:"len"`
		Njs    int    `json:"njs"`
		Njt    int    `json:"njt"`
		Bmb    int    `json:"bmb"`
		Nts    int    `json:"nts"`
		Obs    int    `json:"obs"`
	} `json:"diffs"`
	Key           string      `json:"key"`
	Mapper        string      `json:"mapper"`
	Song          string      `json:"song"`
	Bpm           int         `json:"bpm"`
	DownloadCount int         `json:"downloadCount"`
	UpVotes       int         `json:"upVotes"`
	DownVotes     int         `json:"downVotes"`
	Heat          float64     `json:"heat"`
	Rating        float64     `json:"rating"`
	Automapper    interface{} `json:"automapper"`
	Uploaddate    time.Time   `json:"uploaddate"`
}

var allData map[string]data
var rankedSongs map[string]bool

func main() {
	rankedSongs = make(map[string]bool)
	ranked, err := zip.OpenReader("ranked_all.zip")
	if err != nil {
		panic(err)
	}
	defer ranked.Close()

	for _, file := range ranked.File {
		f, err := file.Open()
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		pl := bplist.NewPlaylist()
		err = json.Unmarshal(data, &pl)
		for _, in := range pl.Songs {
			rankedSongs[in.Hash] = true
		}
		f.Close()
	}

	f, err := os.Open("v2-all.json")
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	allData = make(map[string]data)
	json.Unmarshal(file, &allData)
	log.Println("read file")

	c := len(allData)
	cMutex := sync.Mutex{}
	hashCh := make(chan string, c)

	wg := sync.WaitGroup{}
	for hash := range allData {
		hashCh <- hash
	}

	wg.Add(c)
	log.Println(c, "to process")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			fmt.Printf("%d/%d\n", len(allData)-c, len(allData))
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 1):
				continue
			}
		}
	}()

	for i := 0; i < 2; i++ {
		go func() {
			for {
				hash := <-hashCh
				if _, ok := rankedSongs[hash]; ok {
					getData(hash)
				}
				wg.Done()
				cMutex.Lock()
				c--
				cMutex.Unlock()
			}
		}()
	}

	wg.Wait()
	cancel()
	log.Println("done")

	// write back to disk
	nf, err := os.Create("v2-all-fixed-sync.json")
	if err != nil {
		panic(err)
	}
	defer nf.Close()
	json.NewEncoder(nf).Encode(allData)
}

func getData(hash string) {
	for i := range allData[hash].Diffs {
		allData[hash].Diffs[i].Pp = getPP(hash, allData[hash].Diffs[i].Diff)
		if allData[hash].Diffs[i].Pp != "0" {
			log.Println(hash, allData[hash].Diffs[i].Diff, allData[hash].Diffs[i].Pp)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func getPP(hash, diff string) string {
	// TODO: fetch game type
	res, err := http.Get(fmt.Sprintf("https://api.beatsaver.com/scores/%s/1?difficulty=%d&gamemode=0", hash, difficultyToInt(diff)))
	if err != nil {
		log.Println(err)
		return "0"
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	data := score{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(string(body))
		log.Println(err)
		// probably hit a rate limit
		time.Sleep(time.Second)
		return getPP(hash, diff)
	}
	if !data.Ranked {
		return "0"
	}
	if len(data.Scores) == 0 {
		return "0"
	}

	return fmt.Sprintf("%d", int(math.Round(data.Scores[0].Pp)))
}

func difficultyToInt(diff string) int {
	switch diff {
	case "Easy":
		return 1
	case "Normal":
		return 3
	case "Hard":
		return 5
	case "Expert":
		return 7
	case "ExpertPlus", "Expert+":
		return 9
	}
	return 0
}
