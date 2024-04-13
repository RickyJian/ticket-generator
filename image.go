package main

import (
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
	lineSpacing  = 1.5
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
	ticketInfo := info.string()
	bottomWidth, bottomHeight := canvas.MeasureMultilineString(ticketInfo, lineSpacing)
	bottomX := (widthFloat - bottomWidth) / 2
	bottomY := position + ((heightFloat - position - bottomHeight) / 2)
	canvas.DrawStringWrapped(ticketInfo, bottomX, bottomY, 0, 0, widthFloat, lineSpacing, gg.AlignLeft)

	if err := canvas.SavePNG(dest); err != nil {
		return err
	}
	return nil
}

const (
	dpi        = 300
	pixelRatio = 2.54
)

func pixelFloat(centimeter float32) float64 {
	return math.Round(float64(centimeter) * dpi / pixelRatio)
}

func pixelInt(centimeter float32) int {
	return int(pixelFloat(centimeter))
}
