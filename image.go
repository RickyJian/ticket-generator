package main

import (
	"errors"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"path"

	"github.com/RickyJian/gg"
	"github.com/disintegration/imaging"
)

const (
	white = "#FFFFFF"
	black = "#000000"
	red   = "#FF0000"
	// below unit is centimeter
	width             = 14.8
	height            = 10
	ticketHeight      = 10
	ticketWidth       = 6.6
	dividerPosition   = 7.5
	psvSpace          = 5
	lineSpacing       = 1.5
	gradientStopCount = 5
	opacity           = 255
)

func process(fontFamily string, ts []*tickets) error {
	// drawing canvas
	canvas := newDrawer(height, width)
	// drawing tickets
	for i, t := range ts {
		var x float64
		if i%ticketCount > 0 {
			x = pixelFloat(width) - canvas.ticketWidthPxFloat
		}
		if err := canvas.drawTicket(fontFamily, x, t); err != nil {
			return err
		}
	}
	return canvas.save()
}

type drawer struct {
	*gg.Context
	ticketWidthPx       int
	ticketHeightPx      int
	ticketWidthPxFloat  float64
	ticketHeightPxFloat float64
}

func newDrawer(height, width float32) *drawer {
	canvas := gg.NewContext(pixelInt(width), pixelInt(height))
	canvas.DrawRectangle(0, 0, pixelFloat(width), pixelFloat(height))
	canvas.SetHexColor(white)
	canvas.Fill()
	return &drawer{
		Context:             canvas,
		ticketWidthPx:       pixelInt(ticketWidth),
		ticketHeightPx:      pixelInt(ticketHeight),
		ticketWidthPxFloat:  pixelFloat(ticketWidth),
		ticketHeightPxFloat: pixelFloat(ticketHeight),
	}
}

func (d *drawer) drawTicket(fontFamily string, base float64, t *tickets) error {
	if d == nil {
		return nil
	}

	// drawing background image
	if t.Background != "" {
		background, err := gg.LoadImage(path.Join(src, t.Background))
		if err != nil {
			log.Fatalf("failed to read background image. err: %v\n", err)
		}
		fillImage := imaging.Fill(background, d.ticketWidthPx, d.ticketHeightPx, imaging.Center, imaging.Lanczos)
		d.DrawImage(fillImage, int(base), 0)
		d.Fill()

		// add cover
		coverY := d.ticketHeightPxFloat - d.ticketHeightPxFloat/2
		grad := gg.NewLinearGradient(base, coverY, base, d.ticketHeightPxFloat)
		for i := 0; i < gradientStopCount; i++ {
			offset := float64(i) / float64(gradientStopCount)
			var alpha uint8
			if transCount := gradientStopCount - 2; i < transCount {
				alpha = uint8(opacity / transCount * i)
			} else {
				alpha = opacity
			}
			gradientColor := color.NRGBA{R: 255, G: 255, B: 255, A: alpha}
			grad.AddColorStop(offset, gradientColor)
		}
		d.SetFillStyle(grad)
		d.DrawRectangle(base, coverY, d.ticketWidthPxFloat, d.ticketHeightPxFloat)
		d.Fill()
	}

	position := pixelFloat(dividerPosition)
	// drawing divider position
	d.DrawLine(base, position, base+d.ticketWidthPxFloat, position)
	d.SetHexColor(black)
	d.SetLineWidth(1)
	d.SetDash(3, 5)
	d.Stroke()
	d.DrawLine(base, position, base+d.ticketWidthPxFloat, position)
	d.SetHexColor(black)
	d.SetLineWidth(1)
	d.SetDash(3, 5)
	d.Stroke()

	// drawing text
	// 免責訊息
	notifyFontSize := 20.0
	if err := d.LoadFontFace(path.Join(src, fontFamily), notifyFontSize); err != nil {
		return err
	}
	d.SetHexColor(red)
	notifyText := "開映前二十分鐘，恕不退換。"
	notifyWidth, _ := d.MeasureString(notifyText)
	bottomTopX, bottomTopY := base+((d.ticketWidthPxFloat-notifyWidth)/2), position-psvSpace
	d.DrawStringAnchored(notifyText, bottomTopX, bottomTopY, 0, 0)

	notifyText2 := fmt.Sprintf("限%s使用，隔日無效。", t.Cinema.Name)
	notifyWidth2, _ := d.MeasureString(notifyText2)
	bottomRightX2, bottomRightY2 := base+(d.ticketWidthPxFloat-notifyWidth2), d.ticketHeightPxFloat-psvSpace
	d.DrawStringAnchored(notifyText2, bottomRightX2, bottomRightY2, 0, 0)

	// 電影票訊息
	infoFontSize := 32.0
	if err := d.LoadFontFace(path.Join(src, fontFamily), infoFontSize); err != nil {
		return err
	}
	d.SetHexColor(black)
	ticketInfo := t.string()
	bottomWidth, bottomHeight := d.MeasureMultilineString(ticketInfo, lineSpacing)
	bottomX := base + ((d.ticketWidthPxFloat - bottomWidth) / 2)
	bottomY := position + ((d.ticketHeightPxFloat - position - bottomHeight) / 2)
	d.DrawStringWrapped(ticketInfo, bottomX, bottomY, 0, 0, d.ticketWidthPxFloat, lineSpacing, gg.AlignLeft)
	return nil
}

func (d *drawer) save() error {
	if d == nil {
		return errors.New("empty drawer")
	}
	return d.SavePNG(dest)
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
