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

	"github.com/schollz/progressbar/v3"
)

const (
	dataName    = "data.json"
	ticketCount = 2
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
	} else if _, err := os.Stat(src); os.IsNotExist(err) {
		log.Fatal("src directory not found")
	} else if _, err := os.Stat(dest); os.IsNotExist(err) {
		log.Fatal("desc directory not found")
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

	bar := progressbar.NewOptions(len(info.Tickets),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetItsString("pieces"),
		progressbar.OptionSetDescription("Generating tickets..."),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetTheme(
			progressbar.Theme{
				Saucer:        "=",
				SaucerHead:    ">",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]"}))
	for len(info.Tickets) > 0 {
		size := ticketCount
		if size > len(info.Tickets) {
			size = len(info.Tickets)
		}
		ts := info.Tickets[:size]
		if err := process(info.FontFamily, ts); err != nil {
			log.Printf("failed to process image. err: %v\n", err)
			continue
		}
		if err := bar.Add(size); err != nil {
			log.Fatal(err)
		}
		info.Tickets = info.Tickets[size:]
	}
}

type data struct {
	FontFamily string `json:"font_family"`
	// ticket information
	Tickets []*tickets
}

func (d *data) validate() error {
	if d.FontFamily == "" {
		return errors.New("font family cannot be empty")
	}
	for _, t := range d.Tickets {
		if t.Cinema.Name == "" {
			return errors.New("cinema name cannot be empty")
		} else if t.Movie.Name == "" {
			return errors.New("movie name cannot be empty")
		} else if t.Movie.Name == "" {
			return errors.New("movie eng name cannot be empty")
		} else if t.Movie.Time == "" {
			return errors.New("movie time cannot be empty")
		} else if t.Ticket.Room == "" {
			return errors.New("ticket room cannot be empty")
		} else if t.Ticket.Seat == "" {
			return errors.New("ticket seat cannot be empty")
		} else if t.Ticket.Type == "" {
			return errors.New("ticket type cannot be empty")
		} else if t.Ticket.Price <= 0 {
			return errors.New("ticket price must be greater than 0")
		}
	}
	return nil
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

type tickets struct {
	Background string  `json:"background"`
	Cinema     *cinema `json:"cinema"`
	Movie      *movie  `json:"movie"`
	Ticket     *ticket `json:"ticket"`
}

func (t *tickets) string() string {
	var builder strings.Builder
	// 影城
	builder.WriteString("影城：")
	builder.WriteString(t.Cinema.Name)
	builder.WriteString("\n")

	// 片名
	builder.WriteString("片名：")
	builder.WriteString(t.Movie.Name)
	if t.Movie.EngName != "" {
		builder.WriteString("（")
		builder.WriteString(t.Movie.EngName)
		builder.WriteString("）")
	}
	builder.WriteString("\n")

	// 放映時間
	builder.WriteString("時間：")
	builder.WriteString(t.Movie.Time)
	if t.Ticket.SalesTime != "" {
		builder.WriteString("（售出：")
		builder.WriteString(t.Ticket.SalesTime)
		builder.WriteString("）")
	}
	builder.WriteString("\n")

	// 影廳
	builder.WriteString("影廳：")
	builder.WriteString(t.Ticket.Room)
	builder.WriteString("\n")

	// 座位
	builder.WriteString("座位：")
	builder.WriteString(t.Ticket.Seat)
	builder.WriteString("\n")

	// 票價
	builder.WriteString("票價：")
	builder.WriteString(strconv.Itoa(t.Ticket.Price))
	builder.WriteString(" 元")
	builder.WriteString("（")
	builder.WriteString(t.Ticket.Type)
	builder.WriteString("）")
	return builder.String()
}
