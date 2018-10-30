package main

import (
	"image/color"
	"fmt"
    "math/rand"
    "time"
	"image/png"
	"os"
	"image"
	"log"

	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/anthonynsimon/bild/noise"
	"github.com/anthonynsimon/bild/effect"
	//"github.com/andybons/gogif"
)

//"github.com/rtropisz/monroe/monroe"
func random(min, max int) int {
	return min + rand.Intn(max-min)
}

const inputPath = "resources/CMYK_split/"
const outPath = "img/"
const shift = 10
// Should be 3-5
const noiseSize = 3

func main() {

	C := read(inputPath+"cyan")
	Y := read(inputPath+"yellow")
	M := read(inputPath+"magenta")
	mask := read(inputPath+"mask_temp")
	K := read(inputPath+"black")
	noise := effect.Invert(makeSomeNoise(C.Bounds().Max.X))
	noiseAsAlpha := rgbaToAlpha(blend.Subtract(effect.Invert(mask), noise))
	save("new_mask", noiseAsAlpha)

	fC := toRGBA(C)
	fM := toRGBA(M)
	fY := toRGBA(Y)

	r := rand.Intn(3) 
	switch r {
	case 0:
		fmt.Printf("Cyan\n")
		// fC = maskuj(fC, noiseAsAlpha)
		fC = maskRGBA(fC, noiseAsAlpha)
		save("fC",fC)
	case 1:
		fmt.Printf("Magenta\n")
		fM = maskRGBA(fM, noiseAsAlpha)
		save("fM",fM)
	case 2:
		fmt.Printf("Yellow\n")
		fY = maskRGBA(fY, noiseAsAlpha)
		save("fY",fY)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	r = rand.Intn(3) 

	switch r {
	case 0:
		fC = shiftImage(fC)
	case 1:
		fM = shiftImage(fM)
	case 2:
		fY = shiftImage(fY)
	}

	C = blend.Multiply(K, fC)
	CM := blend.Multiply(C, fM)
	CMY := blend.Multiply(CM, fY)
	

	save("CK", C)
	save("CMK", CM)
	save("CMYK", CMY)

}

func rgbaToAlpha(img *image.RGBA) *image.Alpha {
	out := image.NewAlpha(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			oldColor := img.RGBAAt(x, y)
            color := color.Alpha{oldColor.G}
			out.Set(x,y,color)
		}
	}
	return out
}

func shiftImage(src *image.RGBA) *image.RGBA {
	x := random(-shift,shift)
	y := random(-shift,shift)
	return transform.Translate(src,x,y)
}

func toRGBA(img image.Image) *image.RGBA {
	out := image.NewRGBA(img.Bounds())
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			oldColor := img.At(x, y)
            color := color.RGBAModel.Convert(oldColor)
			out.Set(x,y,color)
		}
	}
	return out
}

func invert(img *image.Alpha) *image.Alpha {
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.SetAlpha(x,y,color.Alpha{255-img.AlphaAt(x,y).A})
		}
	}
	return img
} 

func makeSomeNoise(size int) *image.RGBA {
	noise := noise.Generate(
		noiseSize,
		noiseSize,
		&noise.Options{NoiseFn: noise.Uniform, Monochrome: true})
	noise = transform.Resize(
		noise,
		size,
		size,
		transform.Lanczos,
	)
	return noise
}

func maskRGBA(img *image.RGBA, mask *image.Alpha) *image.RGBA {
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			alpha := mask.AlphaAt(x,y).A
			color := color.NRGBA{
				img.RGBAAt(x,y).R,
				img.RGBAAt(x,y).G,
				img.RGBAAt(x,y).B,
				alpha,
			}
			img.Set(x,y,color)
		}
	}
	return img
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
}