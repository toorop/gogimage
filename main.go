package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"net/http"

	"github.com/fogleman/gg"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	width  = 1200
	height = 628
)

func main() {

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello world")
	})

	// og image generatro
	e.GET("/img/:title/:footer", ctrlOgImage)
	e.GET("/img/:title", ctrlOgImage)
	e.Logger.Fatal(e.Start("127.0.0.1:9999"))
}

// ctrlOgImage controler which return generated image
func ctrlOgImage(c echo.Context) error {

	title, err := base64.StdEncoding.DecodeString(c.Param("title"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	footer := []byte{}
	if len(c.Param("footer")) != 0 {
		footer, err = base64.StdEncoding.DecodeString(c.Param("footer"))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	b, err := generateOGIm(string(title), string(footer))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	blob := b.Bytes()
	c.Response().Header().Add("content-length", fmt.Sprintf("%d", len(blob)))
	return c.Blob(http.StatusOK, "image/png", b.Bytes())

}

// generateOGIm Open Graph Image generator
func generateOGIm(title, footer string) (b bytes.Buffer, err error) {
	// init gg
	ggCtx := gg.NewContext(width, height)
	bgImage, err := gg.LoadImage("bg.png")
	if err != nil {
		return b, err
	}
	ggCtx.DrawImage(bgImage, 0, 0)

	// Font
	fontPath := "Roboto-Regular.ttf"

	// footer
	textColor := color.RGBA{255, 204, 0, 255}
	ggCtx.SetColor(textColor)
	if err := ggCtx.LoadFontFace(fontPath, 30); err != nil {
		return b, err
	}
	textRightMargin := 60.0
	maxWidth := float64(ggCtx.Width()) - textRightMargin - textRightMargin
	x := textRightMargin
	y := float64(ggCtx.Height()) - 75
	ggCtx.DrawStringWrapped(footer, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	// titre
	textColor = color.RGBA{255, 255, 255, 255}
	if err := ggCtx.LoadFontFace(fontPath, 100); err != nil {
		return b, err
	}
	x = textRightMargin
	y = 60.0
	ggCtx.SetColor(textColor)
	ggCtx.DrawStringWrapped(title, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	err = ggCtx.EncodePNG(&b)
	return b, err
}
