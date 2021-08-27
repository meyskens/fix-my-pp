package main

import "time"

type scrapeEntry struct {
	Key             string `json:"Key"`
	Hash            string `json:"Hash"`
	SongName        string `json:"SongName"`
	SongSubName     string `json:"SongSubName"`
	SongAuthorName  string `json:"SongAuthorName"`
	LevelAuthorName string `json:"LevelAuthorName"`
	Diffs           []struct {
		Diff             string        `json:"Diff"`
		Char             string        `json:"Char"`
		Stars            float64       `json:"Stars"`
		Ranked           bool          `json:"Ranked"`
		RankedUpdateTime time.Time     `json:"RankedUpdateTime"`
		Bombs            int           `json:"Bombs"`
		Notes            int           `json:"Notes"`
		Obstacles        int           `json:"Obstacles"`
		Njs              int           `json:"Njs"`
		NjsOffset        int           `json:"NjsOffset"`
		Requirements     []interface{} `json:"Requirements"`
	} `json:"Diffs"`
	Chars     []string  `json:"Chars"`
	Uploaded  time.Time `json:"Uploaded"`
	Uploader  string    `json:"Uploader"`
	Bpm       int       `json:"Bpm"`
	Downloads int       `json:"Downloads"`
	Upvotes   int       `json:"Upvotes"`
	Downvotes int       `json:"Downvotes"`
	Duration  int       `json:"Duration"`
}
