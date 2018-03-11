package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"log"
)

var img *ebiten.Image

func update(image *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2,3)
	op.GeoM.Rotate(0.234)           //random non 45-degree angle
	op.GeoM.Translate(100, 100)     //just for a fine pic, not important
	GeomR:=op.GeoM
	GeomR.Invert()
	for x:=-10.0;x<10.0;x++{
		for x:=-10.0;x<10.0;x++{

		}

	}
	op.Filter = ebiten.FilterLinear //non DefaultFilter here needed
	image.DrawImage(img, op)
	return nil
}


func main() {
	//order is important: first load drawable image
	filter := ebiten.FilterDefault                                           //any, doesn't matter

	//then load another image
	//no matter png or jpg
	//as you see we do not do anything with it later
	//ebitenutil.NewImageFromFile("noise3.png", filter) without variable do the same
	img2, _, _ := ebitenutil.NewImageFromFile("small.png", filter)
	_ = img2
	img, _, _ = ebitenutil.NewImageFromFile("ship.png", ebiten.FilterLinear) //non DefaultFilter here needed

	GeoM:=ebiten.GeoM{}
	GeoM.Scale(2,3)
	GeoM.Rotate(0.234)           //random non 45-degree angle
	GeoM.Translate(100, 100)     //just for a fine pic, not important
	GeomR:=GeoM
	GeomR.Invert()
	for x:=-10.0;x<10.0;x++{
		for y:=-10.0;y<10.0;y++{
			nx,ny:=GeoM.Apply(x,y)
			mx,my:=GeomR.Apply(nx,ny)
			if math.Abs(mx-x)>0.001 || math.Abs(my-y)>0.1 {
				log.Println(x,y,nx,ny,mx,my)
			}
		}

	}

	ebiten.Run(update, 800, 800, 1, "minitest")
}
