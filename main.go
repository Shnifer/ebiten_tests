package main

import (
	"fmt"
	"github.com/Shnifer/ebiten_tests/graph"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"log"
	"time"
	"strconv"
)

const winW = 800
const winH = 600

const showtimeeach = 300

var frameN int
var sprite *graph.Sprite
var dur, totdur int64
var startPeriod time.Time
var totPeriod time.Duration
var face font.Face
var text *graph.Text
var tsprite *graph.Sprite
var testimg *ebiten.Image
var cam *graph.Camera

const camSpeed = 100

var last time.Time

var NoiseActive [3]bool
var NoiseSprite [3]*graph.Sprite
var NoisePlayer [3]*audio.Player

const NoiseSize = 128

var audioContext *audio.Context

func turnOn(player *audio.Player){
	if !player.IsPlaying() {
		player.Rewind()
		player.Play()
	}
	player.SetVolume(1)
}

func turnOff(player *audio.Player){
	player.SetVolume(0)
}

func turnAudio(n int, active bool){
	if active {
		turnOn(NoisePlayer[n])
	} else {
		turnOff(NoisePlayer[n])
	}
}

func mainLoop(image *ebiten.Image) error {

	//Проверка на выход
	dt := time.Since(last).Seconds()
	last = time.Now()
	start := time.Now()
	if frameN%showtimeeach == 0 {
		totdur = dur
		dur = 0
		totPeriod = time.Since(startPeriod)
		startPeriod = start
	}
	frameN++
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("Quit")
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		cam.Pos.Y += camSpeed * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		cam.Pos.Y -= camSpeed * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		cam.Pos.X += camSpeed * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		cam.Pos.X -= camSpeed * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		cam.AngleDeg += 20 * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		cam.AngleDeg -= 20 * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		cam.Scale *= (1 + 0.5*dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		cam.Scale /= (1 + 0.5*dt)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		NoiseActive[0] = !NoiseActive[0]
		turnAudio(0,NoiseActive[0])
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		NoiseActive[1] = !NoiseActive[1]
		turnAudio(1,NoiseActive[1])
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		NoiseActive[2] = !NoiseActive[2]
		turnAudio(2,NoiseActive[2])
	}

	cam.Recalc()

	if ebiten.IsRunningSlowly() {
		log.Println("IsRunningSlowly")
		return nil
	}
	//очистка
	image.Fill(colornames.Black)
	//картинка
	sprite.SetAng(float64(frameN/30)*0 + 27)
	img, op := sprite.ImageOp()
	_ = img
	if !NoiseActive[0] {
		op.Filter = ebiten.FilterLinear
	} else {
		op.Filter = ebiten.FilterDefault
	}
	image.DrawImage(testimg, op)

	//текст
	tim,top:=tsprite.ImageOp()
	image.DrawImage(tim,top)
	//text.Draw(image)

	//Шумы
	for i := 0; i < 3; i++ {
		if !NoiseActive[i] {
			continue
		}
		log.Println(i)
		img, op := NoiseSprite[i].ImageOp()
		image.DrawImage(img, op)
	}

	//ФПС сообщение
	fps := ebiten.CurrentFPS()
	msg := fmt.Sprintf("FPS: %v\nDur:%v micS\nPeriod = %.2f", fps, totdur/1000/showtimeeach, totPeriod.Seconds())
	ebitenutil.DebugPrint(image, msg)
	dur += time.Since(start).Nanoseconds()
	return nil
}

func main() {

	testimg, _, _ = ebitenutil.NewImageFromFile("ship.png", ebiten.FilterLinear)

	//Камера
	cam = new(graph.Camera)
	cam.Scale = 1
	cam.Center = graph.Point{winW / 2, winH / 2}

	//Корабль
	tex, err := graph.GetTex("ship.png", ebiten.FilterLinear, 0, 0)
	if err != nil {
		panic(err)
	}
	sprite = graph.NewSprite(tex, cam, false)
	sprite.SetColor(colornames.Mediumaquamarine)
	sprite.SetSize(400, 400)
	sprite.SetPos(graph.Point{0, 0})

	last = time.Now()

	//текст
	face, err = graph.GetFace("phantom.ttf", 20.0)
	if err != nil {
		panic(err)
	}
	text = graph.NewText("MyText", face, colornames.Darkorchid)
	ttex := graph.TexFromImage(text.Image(), ebiten.FilterLinear, 0, 0)
	tsprite = graph.NewSprite(ttex, cam, true)
	text.SetPosPivot(cam.Center, graph.Center())

	//Шумы
	audioContext,err=audio.NewContext(48000)

	filter := ebiten.FilterDefault
	NoiseSprite[0], _ = graph.NewSpriteFromFile("noise1.jpg", filter, 0, 0, nil, false)
	NoiseSprite[1], _ = graph.NewSpriteFromFile("noise2.jpg", filter, 0, 0, nil, false)
	NoiseSprite[2], _ = graph.NewSpriteFromFile("noise3.png", filter, 0, 0, nil, false)
	for i := 0; i < 3; i++ {
		NoiseSprite[i].SetPos(graph.Point{winW - NoiseSize, NoiseSize})
		NoiseSprite[i].SetAlpha(0.75)
		f,err:=ebitenutil.OpenFile("noise"+strconv.Itoa(i+1)+".wav")
		if err!=nil{
			panic(err)
		}
		d,err:=wav.Decode(audioContext,f)
		if err!=nil{
			panic(err)
		}
		NoisePlayer[i],err=audio.NewPlayer(audioContext,d)
		if err!=nil{
			panic(err)
		}
	}
	NoiseSprite[0].SetColor(colornames.Red)
	NoiseSprite[1].SetColor(colornames.Green)
	NoiseSprite[2].SetColor(colornames.Blue)


	//audio
	ebiten.SetFullscreen(false)
	ebiten.SetRunnableInBackground(true)
	if err := ebiten.Run(mainLoop, winW, winH, 1, "Ebiten Test"); err != nil {
		log.Fatal(err)
	}
}
