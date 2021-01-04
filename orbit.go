package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type circle struct {
	x   float64
	y   float64
	r   float64
	clr color.Color
}

func NewCircle(x, y, radius float64, clr color.Color) *circle {
	return &circle{
		x:   x,
		y:   y,
		r:   radius,
		clr: clr,
	}
}

func (c *circle) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *circle) Bounds() image.Rectangle {
	minX, minY := int(c.x-c.r), int(c.y-c.r)
	maxX, maxY := int(c.x+c.r), int(c.y+c.r)
	return image.Rectangle{
		Min: image.Pt(minX, minY),
		Max: image.Pt(maxX, maxY),
	}
}

func (c *circle) At(x, y int) color.Color {
	xx, yy := float64(x)-c.x, float64(y)-c.y
	if xx*xx+yy*yy < c.r*c.r {
		return c.clr
	}
	return color.Alpha{0}
}

type Game struct {
	sprites []*circle
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, sprite := range g.sprites {
		img := ebiten.NewImageFromImage(sprite)
		geom := ebiten.GeoM{}
		geom.Translate(sprite.x, sprite.y)
		screen.DrawImage(img, &ebiten.DrawImageOptions{GeoM: geom})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		sprites: []*circle{
			NewCircle(10, 10, 5, color.RGBA{255, 0, 0, 255}),
		},
	}

	ebiten.SetWindowSize(400, 300)
	ebiten.SetWindowTitle("cargo")

	log.Fatal(ebiten.RunGame(game))
}
