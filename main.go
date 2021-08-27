package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"
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

func main() {
	f, err := os.Open("v2-all.json")
	if err != nil {
		panic(err)
	}

	file, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	allData := map[string]data{}
	json.Unmarshal(file, &allData)

	c := 0
	for hash, level := range allData {
		for _, diff := range level.Diffs {
			diff.Pp = getPP(hash, diff.Diff)
			if diff.Pp != "0" {
				log.Println(hash, diff.Diff, diff.Pp)
			}
			time.Sleep(5 * time.Millisecond)
		}
		c++
		fmt.Printf("\r%d/%d", c, len(allData))
	}

	// write back to disk
	nf, err := os.Create("v2-all-fixed.json")
	if err != nil {
		panic(err)
	}
	defer nf.Close()
	json.NewEncoder(nf).Encode(allData)
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
	json.Unmarshal(body, &data)
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
	case "ExpertPlus":
	case "Expert+":
		return 9
	}
	return 0
}
