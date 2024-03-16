package main

import (
	"encoding/json"
	"log"
	"os"
)

const (
	path = "data/data.json"
)

func main() {
	bs, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read data.json. err: %v\n", err)
	}
	info := &data{}
	if err := json.Unmarshal(bs, &info); err != nil {
		log.Fatalf("failed to unmarshal data. err: %v\n", err)
	}
}

type data struct {
	Width           float32 `json:"width"`
	Length          float32 `json:"length"`
	DividerPosition float32 `json:"divider_position"`
	Cinema          *cinema `json:"cinema"`
	Movie           *movie  `json:"movie"`
	Ticket          *ticket `json:"ticket"`
}

type cinema struct {
	Name    string `json:"name"`
	EngName string `json:"eng_name"`
}

type movie struct {
	Name    string `json:"name"`
	EngName string `json:"eng_name"`
	Time    string `json:"time"`
	Room    string `json:"room"`
}

type ticket struct {
	Seat      string `json:"seat"`
	Type      string `json:"type"`
	Price     int    `json:"price"`
	SalesTime string `json:"sales_time"`
}
