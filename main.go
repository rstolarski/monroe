package main

import (
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
	//"github.com/andybons/gogif"
)

//"github.com/rtropisz/monroe/monroe"
func random(min, max int) int {
	return min + rand.Intn(max-min)
}

const inputPath = "resources/CMYK_split/"
const outPath = "img/"
const i = 5

func main() {

	C := read(inputPath+"cyan")
	Y := read(inputPath+"yellow")
	M := read(inputPath+"magenta")
	//K := read(inputPath+"black")

	rand.Seed(time.Now().UTC().UnixNano())
	x := random(-i,i)
	y := random(-i,i)
	C = transform.Translate(C,x,y)
	x = random(-i,i)
	y = random(-i,i)
	Y = transform.Translate(Y,x,y)
	x = random(-i,i)
	y = random(-i,i)
	M = transform.Translate(M,x,y)
	
	// Don't move Key
	// K = transform.Translate(K,-5,6)

	CM := blend.Multiply(C, M)
	CMY := blend.Multiply(CM, Y)
	//CMYK := blend.Multiply(CMY, K)

	save("C", C)
	save("CM", CM)
	save("CMY", CMY)
	//save("CMYK", CMYK)
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