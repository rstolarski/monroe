package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/anthonynsimon/bild/adjust"

	"github.com/disintegration/imaging"

	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

func random(min, max int) int {
	return min + rand.Intn(max-min)
}

const inputPath = "resources/CMYK_split/"
const outPath = "img/"

const loops = 6

const fps = 25

func main() {
	//b := bytes.NewBuffer()
	bgColor := color.RGBA{249, 92, 137, 255} // magenta
	bgColor = color.RGBA{82, 145, 188, 255}  //cyan

	var f []*image.RGBA
	for i := 1; i < loops; i += 2 {
		f = append(f, mainLoop(i, bgColor)...)
	}
	numberOfFrames := len(f)
	log.Printf("Number of frames: %v\n", numberOfFrames)
	log.Printf("Exporting frames\n")
	for i, img := range f {
		save(fmt.Sprintf("%06d", i), img)
	}
}

// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

func selectMask(outShift int, pt image.Point) *image.RGBA {
	rX := random(0, outShift)
	rand.Seed(time.Now().UTC().UnixNano())
	rY := random(0, outShift)

	anchor := image.Point{(pt.X / outShift) * rX, (pt.Y / outShift) * rY}

	mask := read(inputPath + "mask")
	dst := imaging.New(pt.X, pt.Y, color.RGBA{0, 0, 0, 0})
	maksDst := imaging.Resize(mask, pt.X/outShift, 0, imaging.Lanczos)
	maksDst = imaging.Paste(dst, maksDst, anchor)
	log.Printf("x,y: %v,%v", anchor.X, anchor.Y)
	return toRGBA(maksDst)
}

func changeHueToRandom(img *image.RGBA) *image.RGBA {
	c := random(0, 36)
	return adjust.Hue(img, c*10)
}

func mainLoop(outShift int, bgColor color.RGBA) []*image.RGBA {
	C := read(inputPath + "cyan")
	Y := read(inputPath + "yellow")
	M := read(inputPath + "magenta")
	K := read(inputPath + "black")

	//m := read(inputPath + "mask").Bounds().Max
	//_ = m

	var f []*image.RGBA
	var tempF []*image.NRGBA

	ptX := C.Bounds().Max.X
	ptY := C.Bounds().Max.Y
	p := image.Point{ptX / outShift, ptY / outShift}
	dst := imaging.New(ptX, ptY, bgColor)
	for x := 0; x < outShift; x++ {
		for y := 0; y < outShift; y++ {
			dstImg := convertAndShiftAllImages(K, C, M, Y, ptX/outShift, 0, 1)
			dst, tempF = combineAllImages(f, dst, dstImg, image.Point{p.X * x, p.Y * y})
			for _, i := range tempF {
				img := setBG(bgColor, toRGBA(i))
				//for n := 0; n < fps; n++ {
				f = append(f, img)
				//}
			}
		}
	}
	final := toRGBA(dst)
	//for i := 0; i < fps; i++ {
	f = append(f, setBG(bgColor, final))
	//}

	for i := 0; i < outShift*outShift; i++ {
		mask := selectMask(outShift, image.Point{C.Bounds().Max.X, C.Bounds().Max.Y})
		mask = changeHueToRandom(maskRGBA(final, mask))
		final = combineMaskedWithFinalFrame(final, mask)
		final = setBG(bgColor, final)
		//for n := 0; n < fps; n++ {
		f = append(f, final)
		//}
	}

	//for i := 0; i < fps; i++ {
	f = append(f, final)
	//	}
	return f
}

func combineMaskedWithFinalFrame(f, m *image.RGBA) *image.RGBA {
	r := image.NewRGBA(f.Bounds())
	for y := 0; y < f.Bounds().Max.Y; y++ {
		for x := 0; x < f.Bounds().Max.X; x++ {
			mC := m.RGBAAt(x, y)
			fC := f.RGBAAt(x, y)
			if mC.A != 0 && fC.A != 0 {
				r.SetRGBA(x, y, mC)
			} else {
				r.SetRGBA(x, y, fC)
			}
		}
	}
	return r
}

func combineAllImages(
	f []*image.RGBA,
	dst *image.NRGBA,
	dstImg []*image.RGBA,
	p image.Point) (*image.NRGBA, []*image.NRGBA) {

	var frames []*image.NRGBA
	for _, img := range dstImg {
		dst = imaging.Paste(dst, img, p)
		frames = append(frames, dst)
	}
	return dst, frames
}

func convertAndShiftAllImages(
	K image.Image,
	C image.Image,
	M image.Image,
	Y image.Image,
	s, omin, omax int) []*image.RGBA {

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
		fC = shiftImage(fC, "cyan", omin, omax)
	case 1:
		fM = shiftImage(fM, "magenta", omin, omax)
	case 2:
		fY = shiftImage(fY, "yellow", omin, omax)
	}

	dstC := blend.Multiply(fK, fM)
	dstCM := blend.Multiply(dstC, fC)
	dstCMY := blend.Multiply(dstCM, fY)

	return []*image.RGBA{
		//fK,
		dstC,
		dstCM,
		dstCMY,
	}
}

func rgbaToAlpha(img *image.RGBA) *image.Alpha {
	out := image.NewAlpha(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			oldColor := img.RGBAAt(x, y)
			color := color.Alpha{oldColor.G}
			out.Set(x, y, color)
		}
	}
	return out
}

func setBG(c color.RGBA, img *image.RGBA) *image.RGBA {
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			if img.RGBAAt(x, y).A == 0 {
				img.SetRGBA(x, y, c)
			}
		}
	}
	return img
}

func shiftImage(src *image.RGBA, l string, omin, omax int) *image.RGBA {
	x := random(omin, omax)
	y := random(omin, omax)
	rand.Seed(time.Now().UTC().UnixNano())
	r := rand.Intn(2)
	if r > 0 {
		x = -x
		r := rand.Intn(2)
		if r > 0 {
			y = -y
		}
	}
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

func toAlpha(img image.Image) *image.Alpha {
	out := image.NewAlpha(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			oldColor := img.At(x, y)
			color := color.AlphaModel.Convert(oldColor)
			out.Set(x, y, color)
		}
	}
	return out
}

func invert(img *image.RGBA) *image.RGBA {
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{
				255 - img.RGBAAt(x, y).R,
				255 - img.RGBAAt(x, y).G,
				255 - img.RGBAAt(x, y).B,
				255,
			})
		}
	}
	return img
}

func maskRGBA(img, mask *image.RGBA) *image.RGBA {
	out := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			alpha := mask.RGBAAt(x, y).A
			color := color.NRGBA{
				img.RGBAAt(x, y).R,
				img.RGBAAt(x, y).G,
				img.RGBAAt(x, y).B,
				alpha,
			}
			out.Set(x, y, color)
		}
	}
	return out
}

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
		log.Fatalf("MkdirAll %q: %s", outPath, err)
	}
	f, err := os.Create(outPath + output + ".jpg")
	if err != nil {
		log.Fatal(err)
	}

	if err := jpeg.Encode(f, image, nil); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	//log.Printf("Image was saved\n")
}
