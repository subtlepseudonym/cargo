package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pbnjay/pixfont"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

var (
	useStep     = false
	step        = false
	Orange      = color.RGBA{255, 165, 0, 255}
	DeepSkyBlue = color.RGBA{0, 191, 255, 255}
	FireBrick   = color.RGBA{178, 34, 34, 255}
)

type Body struct {
	r   float64
	clr color.Color

	x  float64
	y  float64
	dx float64
	dy float64

	mass float64
}

func (b *Body) ColorModel() color.Model {
	return color.RGBAModel
}

func (b *Body) Bounds() image.Rectangle {
	return image.Rect(int(b.x-b.r), int(b.y-b.r), int(b.x+b.r), int(b.y+b.r))
}

func (b *Body) At(x, y int) color.Color {
	xx, yy := float64(x-int(b.x))+0.5, float64(y-int(b.y))+0.5 // FIXME
	if xx*xx+yy*yy < b.r*b.r {
		return b.clr
	}
	return color.Alpha{0}
}

type Game struct {
	sun     *Body
	planets []*Body
}

func (g *Game) Update() error {
	if useStep {
		if !ebiten.IsKeyPressed(ebiten.KeySpace) {
			if step {
				step = false
			}
			return nil
		} else if step {
			return nil
		}
		step = true
	}

	for _, planet := range g.planets {
		gravityX, gravityY := gravity(g.sun, planet)
		planet.dx = planet.dx + gravityX
		planet.dy = planet.dy + gravityY

		planet.x = planet.x + planet.dx
		planet.y = planet.y + planet.dy
	}

	return nil
}

func gravity(sun, planet *Body) (float64, float64) {
	distX := sun.x - planet.x
	distY := sun.y - planet.y
	distSq := math.Pow(distX, 2) + math.Pow(distY, 2)
	dist := math.Sqrt(distSq)

	f := (sun.mass * planet.mass) / distSq
	fx := (distX / dist) * f
	fy := (distY / dist) * f

	return fx, fy
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBody(screen, g.sun)
	for i, planet := range g.planets {
		drawBody(screen, planet)
		pos := fmt.Sprintf("%.2f, %.2f", planet.x, planet.y)
		pixfont.DrawString(screen, 0, i*10, pos, color.White)
	}
}

func drawBody(screen *ebiten.Image, body *Body) {
	img := ebiten.NewImageFromImage(body)
	geom := ebiten.GeoM{}
	geom.Translate(body.x-body.r, body.y-body.r)

	screen.DrawImage(img, &ebiten.DrawImageOptions{
		GeoM:   geom,
		Filter: ebiten.FilterNearest,
	})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	sol := &Body{
		x:    screenWidth / 2,
		y:    screenHeight / 2,
		r:    40,
		clr:  Orange,
		mass: 10,
	}
	earth := &Body{
		x:    screenWidth / 2,
		y:    80,
		r:    10,
		clr:  DeepSkyBlue,
		dx:   -0.28,
		mass: 1,
	}
	mars := &Body{
		x:    screenWidth / 2,
		y:    20,
		r:    8,
		clr:  FireBrick,
		dx:   -0.21,
		mass: 0.8,
	}

	game := &Game{
		sun: sol,
		planets: []*Body{
			earth,
			mars,
		},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("cargo")

	log.Fatal(ebiten.RunGame(game))
}
