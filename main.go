package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	dataName = "data.json"
)

var (
	src  string
	dest string
)

func init() {
	flag.StringVar(&src, "s", "", "input folder(REQUIRED)")
	flag.StringVar(&dest, "d", "", "output folder(REQUIRED)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "ticket-generator is a command line tool to generate fake movie ticket.\n\n")
		fmt.Fprintf(os.Stdout, "Usage:\n\n")
		fmt.Fprintf(os.Stdout, "\tticket-generator <flag> [arguments]\n\n")
		fmt.Fprintf(os.Stdout, "The flags are:\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if src == "" {
		log.Fatalf("flag needs an argument: -s")
	} else if dest == "" {
		log.Fatalf("flag needs an argument: -d ")
	}
}

func main() {
	dataBs, err := os.ReadFile(path.Join(src, dataName))
	if err != nil {
		log.Fatalf("failed to read data.json. err: %v\n", err)
	}
	info := &data{}
	if err := json.Unmarshal(dataBs, &info); err != nil {
		log.Fatalf("failed to unmarshal data. err: %v", err)
	} else if err := info.validate(); err != nil {
		log.Fatal(err)
	}
	if err := process(info); err != nil {
		log.Fatalf("failed to process image. err: %v", err)
	}
}

type data struct {
	Width      float32 `json:"width"`
	Height     float32 `json:"height"`
	Background string  `json:"background"`
	FontFamily string  `json:"font_family"`
	// ticket information
	Cinema *cinema `json:"cinema"`
	Movie  *movie  `json:"movie"`
	Ticket *ticket `json:"ticket"`
}

func (d *data) validate() error {
	if d.Width <= 0 {
		return errors.New("width must be greater than 0")
	} else if d.Height <= 0 {
		return errors.New("height must be greater than 0")
	} else if d.FontFamily == "" {
		return errors.New("font family cannot be empty")
	} else if d.Cinema.Name == "" {
		return errors.New("cinema name cannot be empty")
	} else if d.Movie.Name == "" {
		return errors.New("movie name cannot be empty")
	} else if d.Movie.Name == "" {
		return errors.New("movie eng name cannot be empty")
	} else if d.Movie.Time == "" {
		return errors.New("movie time cannot be empty")
	} else if d.Ticket.Room == "" {
		return errors.New("ticket room cannot be empty")
	} else if d.Ticket.Seat == "" {
		return errors.New("ticket seat cannot be empty")
	} else if d.Ticket.Type == "" {
		return errors.New("ticket type cannot be empty")
	} else if d.Ticket.Price <= 0 {
		return errors.New("ticket price must be greater than 0")
	}
	return nil
}

func (d *data) string() string {
	var builder strings.Builder
	// 影城
	builder.WriteString("影城：")
	builder.WriteString(d.Cinema.Name)
	builder.WriteString("\n")

	// 片名
	builder.WriteString("片名：")
	builder.WriteString(d.Movie.Name)
	builder.WriteString("（")
	builder.WriteString(d.Movie.EngName)
	builder.WriteString("）")
	builder.WriteString("\n")

	// 放映時間
	builder.WriteString("時間：")
	builder.WriteString(d.Movie.Time)
	builder.WriteString("（售出：")
	builder.WriteString(d.Ticket.SalesTime)
	builder.WriteString("）")
	builder.WriteString("\n")

	// 影廳
	builder.WriteString("影廳：")
	builder.WriteString(d.Ticket.Room)
	builder.WriteString("\n")

	// 座位
	builder.WriteString("座位：")
	builder.WriteString(d.Ticket.Room)
	builder.WriteString("\n")

	// 票價
	builder.WriteString("票價：")
	builder.WriteString(strconv.Itoa(d.Ticket.Price))
	builder.WriteString(" 元")
	builder.WriteString("（")
	builder.WriteString(d.Ticket.Type)
	builder.WriteString("）")
	return builder.String()
}

type cinema struct {
	Name string `json:"name"`
}

type movie struct {
	Name    string `json:"name"`
	EngName string `json:"eng_name"`
	Time    string `json:"time"`
}

type ticket struct {
	Room      string `json:"room"`
	Seat      string `json:"seat"`
	Type      string `json:"type"`
	Price     int    `json:"price"`
	SalesTime string `json:"sales_time"`
}
