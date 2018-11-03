package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"

	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

//"github.com/rtropisz/monroe/monroe"
func random(min, max int) int {
	return min + rand.Intn(max-min)
}

const inputPath = "resources/CMYK_split/"
const outPath = "img/"
const outShift = 3

func main() {
	C := read(inputPath + "cyan")
	Y := read(inputPath + "yellow")
	M := read(inputPath + "magenta")
	K := read(inputPath + "black")

	frames := make([]*image.NRGBA, outShift*outShift*4)

	ptX := C.Bounds().Max.X
	ptY := C.Bounds().Max.Y
	p := image.Point{ptX / outShift, ptY / outShift}
	dst := imaging.New(ptX, ptY, color.RGBA{0, 0, 0, 0})
	for x := 0; x < outShift; x++ {
		for y := 0; y < outShift; y++ {
			dstImg := convertAndShiftAllImages(K, C, M, Y, ptX/outShift, (5/outShift)*2)
			newP := image.Point{p.X * x, p.Y * y}
			dst = combineAllImages(frames, dst, dstImg, newP)
		}
	}
	log.Printf("Number of frames: %v\n", len(frames))
}

func combineAllImages(
	f []*image.NRGBA,
	dst *image.NRGBA,
	dstImg [4]*image.RGBA,
	p image.Point) *image.NRGBA {

	for i, img := range dstImg {
		dst = imaging.Paste(dst, img, p)
		f = append(f, dst)
		str := strconv.Itoa(p.X) + "_" + strconv.Itoa(p.Y) + "_" + strconv.Itoa(i)
		save(str, dst)
	}
	return dst
}

func convertAndShiftAllImages(
	K image.Image,
	C image.Image,
	M image.Image,
	Y image.Image,
	s, o int) [4]*image.RGBA {

	C = imaging.Resize(C, s, 0, imaging.Lanczos)
	Y = imaging.Resize(Y, s, 0, imaging.Lanczos)
	M = imaging.Resize(M, s, 0, imaging.Lanczos)
	K = imaging.Resize(K, s, 0, imaging.Lanczos)

	fK := toRGBA(K)
	fC := toRGBA(C)
	fY := toRGBA(Y)
	fM := toRGBA(M)

	rand.Seed(time.Now().UTC().UnixNano())
	r := rand.Intn(3)
	switch r {
	case 0:
		fC = shiftImage(fC, "cyan", o)
	case 1:
		fM = shiftImage(fM, "magenta", o)
	case 2:
		fY = shiftImage(fY, "yellow", o)
	}

	dstK := fK
	dstC := blend.Multiply(fK, fC)
	dstCM := blend.Multiply(dstC, fM)
	dstCMY := blend.Multiply(dstCM, fY)

	return [4]*image.RGBA{dstK, dstC, dstCM, dstCMY}
}

// func rgbaToAlpha(img *image.RGBA) *image.Alpha {
// 	out := image.NewAlpha(img.Bounds())
// 	for y := 0; y < img.Bounds().Max.Y; y++ {
// 		for x := 0; x < img.Bounds().Max.X; x++ {
// 			oldColor := img.RGBAAt(x, y)
// 			color := color.Alpha{oldColor.G}
// 			out.Set(x, y, color)
// 		}
// 	}
// 	return out
// }

func shiftImage(src *image.RGBA, l string, o int) *image.RGBA {
	x := random(-o, o)
	y := random(-o, o)
	log.Printf("%v, shifted by: %v, %v\n", strings.ToUpper(l), x, y)
	return transform.Translate(src, x, y)
}

func toRGBA(img image.Image) *image.RGBA {
	out := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			oldColor := img.At(x, y)
			color := color.RGBAModel.Convert(oldColor)
			out.Set(x, y, color)
		}
	}
	return out
}

func invert(img *image.Alpha) *image.Alpha {
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.SetAlpha(x, y, color.Alpha{255 - img.AlphaAt(x, y).A})
		}
	}
	return img
}

// func maskRGBA(img *image.RGBA, mask *image.Alpha) *image.RGBA {
// 	for y := 0; y < img.Bounds().Max.Y; y++ {
// 		for x := 0; x < img.Bounds().Max.X; x++ {
// 			alpha := mask.AlphaAt(x, y).A
// 			color := color.NRGBA{
// 				img.RGBAAt(x, y).R,
// 				img.RGBAAt(x, y).G,
// 				img.RGBAAt(x, y).B,
// 				alpha,
// 			}
// 			img.Set(x, y, color)
// 		}
// 	}
// 	return img
// }

func read(input string) image.Image {
	img, err := imgio.Open(input + ".png")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	return img
}

func save(output string, image image.Image) {
	err := os.MkdirAll(outPath, 0777)
	if err != nil {
		fmt.Errorf("MkdirAll %q: %s", outPath, err)
	}
	f, err := os.Create(outPath + output + ".png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, image); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Image was saved\n")
}
