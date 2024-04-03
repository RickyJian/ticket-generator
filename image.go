package main

import (
	"fmt"
	_ "image/png"
	"log"
	"math"
	"path"

	"github.com/RickyJian/gg"
	"github.com/disintegration/imaging"
)

const (
	white        = "#FFFFFF"
	black        = "#000000"
	red          = "#FF0000"
	topPsvSpace  = 10
	leftPsvSpace = 5
)

func process(info *data) error {
	widthInt, heightInt := pixelInt(info.Width), pixelInt(info.Height)
	widthFloat, heightFloat := pixelFloat(info.Width), pixelFloat(info.Height)
	canvas := gg.NewContext(widthInt, heightInt)
	canvas.DrawRectangle(0, 0, widthFloat, heightFloat)
	canvas.SetHexColor(white)
	canvas.Fill()

	// drawing background image
	if info.Background != "" {
		background, err := gg.LoadImage(path.Join(src, info.Background))
		if err != nil {
			log.Fatalf("failed to read background image. err: %v\n", err)
		}
		fittedImage := imaging.Fit(background, widthInt, heightInt, imaging.Lanczos)
		canvas.DrawImage(fittedImage, 0, 0)
		canvas.Fill()
	}

	position := pixelFloat(info.DividerPosition)
	// drawing divider position
	if info.DividerPosition > 0 {
		canvas.DrawLine(0, position, widthFloat, position)
		canvas.SetHexColor(black)
		canvas.SetLineWidth(1)
		canvas.SetDash(3, 5)
		canvas.Stroke()
	}
	// drawing text
	// 免責訊息
	notifyFontSize := pixelFloat(info.Height-info.DividerPosition) / 10
	if err := canvas.LoadFontFace(path.Join(src, info.FontFamily), notifyFontSize); err != nil {
		return err
	}
	canvas.SetHexColor(red)
	notifyText := "開映前二十分鐘，恕不退換。"
	_, notifyW := canvas.MeasureString(notifyText)
	canvas.DrawStringAnchored(notifyText, (widthFloat+notifyW)/2, position-topPsvSpace/2, 0.5, 0)
	// 電影票訊息
	fontSize := pixelFloat(info.Height-info.DividerPosition) / 8
	if err := canvas.LoadFontFace(path.Join(src, info.FontFamily), fontSize); err != nil {
		return err
	}
	canvas.SetHexColor(black)
	// 影城
	position += topPsvSpace
	cinemaText := fmt.Sprintf("影城：%s", info.Cinema.Name)
	_, cinemaH := canvas.MeasureString(cinemaText)
	canvas.DrawStringAnchored(cinemaText, leftPsvSpace, position, 0, 1)
	// 片名
	position += topPsvSpace + cinemaH
	movieText := fmt.Sprintf("片名：%s", info.Movie.Name)
	_, movieH := canvas.MeasureString(movieText)
	canvas.DrawStringAnchored(movieText, leftPsvSpace, position, 0, 1)
	moviePosition := position + movieH
	movieTitleW, _ := canvas.MeasureString("片名：")
	// 放映時間
	position += topPsvSpace + movieH*2
	movieTimeText := fmt.Sprintf("時間：%s", info.Movie.Time)
	_, movieTimeH := canvas.MeasureString(movieTimeText)
	canvas.DrawStringAnchored(movieTimeText, leftPsvSpace, position, 0, 1)
	// 影廳
	position += topPsvSpace + movieTimeH
	roomText := fmt.Sprintf("影廳：%s", info.Ticket.Room)
	_, roomH := canvas.MeasureString(roomText)
	canvas.DrawStringAnchored(roomText, leftPsvSpace, position, 0, 1)
	// 座位
	seatText := fmt.Sprintf("座位：%s", info.Ticket.Seat)
	canvas.DrawStringAnchored(seatText, widthFloat/2, position, 0, 1)
	// 票別
	position += topPsvSpace + roomH
	typeText := fmt.Sprintf("票別：%s", info.Ticket.Type)
	_, typePositionH := canvas.MeasureString(typeText)
	canvas.DrawStringAnchored(typeText, leftPsvSpace, position, 0, 1)
	// 票價
	priceText := fmt.Sprintf("票價：%d 元", info.Ticket.Price)
	canvas.DrawStringAnchored(priceText, widthFloat/2, position, 0, 1)
	// 原始片名
	oriMovieFontSize := pixelFloat(info.Height-info.DividerPosition) / 10
	if err := canvas.LoadFontFace(path.Join(src, info.FontFamily), oriMovieFontSize); err != nil {
		return err
	}
	canvas.SetHexColor(black)
	engPsvSpace := leftPsvSpace + 2 + movieTitleW
	canvas.DrawStringAnchored(info.Movie.EngName, engPsvSpace, moviePosition+movieH, 0, 0.5)
	// 售出時間
	if err := canvas.LoadFontFace(path.Join(src, info.FontFamily), notifyFontSize); err != nil {
		return err
	}
	canvas.SetHexColor(red)
	salesText := fmt.Sprintf("售出：%s", info.Ticket.SalesTime)
	_, salesTextW := canvas.MeasureString(salesText)
	canvas.DrawStringAnchored(salesText, widthFloat-salesTextW, heightFloat-typePositionH, 1, 1)

	if err := canvas.SavePNG(dest); err != nil {
		return err
	}
	return nil
}

const (
	dpi   = 100
	ratio = 2.54
)

func pixelFloat(centimeter float32) float64 {
	return math.Round(float64(centimeter) * dpi / ratio)
}

func pixelInt(centimeter float32) int {
	return int(pixelFloat(centimeter))
}
