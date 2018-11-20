package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/anthonynsimon/bild/adjust"

	"github.com/disintegration/imaging"

	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

// TODO Get new mask
// TODO Make sure that all images in sequence are colorized

const inputPath = "resources/CMYK_split/"
const width = 1280
const height = 720
const frames = 25
const generateFullSize = true

func appednInLoop(f []*image.RGBA, isLooping bool) []*image.RGBA {
	var r []*image.RGBA
	if !isLooping {
		return r
	}
	for _, v := range f {
		for i := 0; i < frames; i++ {
			r = append(r, v)
		}
	}
	return r
}

func main() {

	// if _, err := os.Stat("out.mp4"); !os.IsNotExist(err) {
	// 	fmt.Printf("Cleaning previously created movie\n")
	// 	runCommand(exec.Command("rm", "out.mp4"))
	// }

	bgColor := color.RGBA{249, 92, 137, 255} // magenta
	bgColor = color.RGBA{82, 145, 188, 255}  //cyan
	//	var f []*image.RGBA
	var f1 []*image.RGBA
	var f2 []*image.RGBA
	var f3 []*image.RGBA
	var f4 []*image.RGBA
	C := read(inputPath + "cyan")
	pt := image.Point{C.Bounds().Max.X, C.Bounds().Max.Y}
	f1 = append(f1, createFrames(bgColor, pt, 1)...)
	f2 = append(f2, createFrames(bgColor, pt, 3)...)
	f4 = append(f4, createFrames(bgColor, pt, 5)...)
	lastSmall := f2[len(f2)-1]
	lastBig := f4[len(f4)-1]
	f3 = append(f3, zoomOut(lastSmall, lastBig)...)
	col := colorize(f3[len(f3)-1], bgColor, pt, 5)

	fmt.Printf("Number of frames: %v\n", len(f1)+len(f2)+len(f3)+len(col))
	fmt.Printf("Exporting frames\n")

	saveGroup(f1, pt, bgColor, "temp/f1/")
	saveGroup(f2, pt, bgColor, "temp/f2/")
	saveGroup(f3, pt, bgColor, "temp/f3/")
	saveGroup(col, pt, bgColor, "temp/col/")

	fmt.Printf("Saving images: 100.00 \n")
	fmt.Printf("Images were saved\n")

	b := []byte("file 'f1.mp4'\nfile 'f2.mp4'\nfile 'f3.mp4'\nfile 'col.mp4'\n")
	f, err := os.Create("temp/mylist.txt")
	if err != nil {
		panic(err)
	}
	n, err := f.Write(b)
	if err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", n)
	defer f.Close()
}
func saveGroup(f []*image.RGBA, pt image.Point, bgColor color.RGBA, outPath string) {
	bg := toRGBA(imaging.New(width, height, bgColor))
	if generateFullSize {
		for i, img := range f {
			save(outPath, fmt.Sprintf("%06d", i), paste(bg, img, image.Point{(width - pt.X) / 2, (height - pt.Y) / 2}))
			fmt.Printf("Saving images: %.2f \n", float64(i)/float64(len(f))*100)
		}
	} else {
		for i, img := range f {
			save(outPath, fmt.Sprintf("%06d", i), img)
			fmt.Printf("Saving images: %.2f \n", float64(i)/float64(len(f))*100)
		}
	}
}

func runCommand(cmd *exec.Cmd) {
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Output %v\n", out.String())
}

func zoomOut(img *image.RGBA, bg *image.RGBA) []*image.RGBA {
	var f []*image.RGBA
	r := resize(img, 338, 338)
	res := paste(bg, r, image.Point{(560 - 338) / 2, (560 - 338) / 2})
	for i := 560; i < 900; i += 12 {
		res = resize(res, i, i)

		x := (i - 560) / 2
		f = append(f, toRGBA(imaging.Crop(res, image.Rect(x, x, x+560, x+560))))
	}
	return reverse(f)
}

func paste(bg, img *image.RGBA, pos image.Point) *image.RGBA {
	return toRGBA(imaging.Paste(bg, img, pos))
}

func resize(img *image.RGBA, width, height int) *image.RGBA {
	return toRGBA(imaging.Resize(img, width, height, imaging.Lanczos))
}

func fill(img image.Image, width, height int) *image.RGBA {
	dstImage := imaging.Fill(img, width, height, imaging.Center, imaging.Lanczos)
	return toRGBA(dstImage)
}

func reverse(a []*image.RGBA) []*image.RGBA {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
	return a
}

func selectMask(outShift int, anchor, pt image.Point) *image.RGBA {
	mask := read(inputPath + "mask")
	dst := imaging.New(pt.X, pt.Y, color.RGBA{0, 0, 0, 0})
	maksDst := imaging.Resize(mask, pt.X/outShift, 0, imaging.Lanczos)
	maksDst = imaging.Paste(dst, maksDst, anchor)
	return toRGBA(maksDst)
}

func changeHueToRandom(img *image.RGBA) *image.RGBA {
	time.Sleep(time.Millisecond)
	rand.Seed(time.Now().UTC().UnixNano())
	c := random(0, 18)
	return adjust.Hue(img, c*20)
}

func remove(s []int, r int) []int {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func createFrames(bgColor color.RGBA, pt image.Point, outShift int) []*image.RGBA {
	C := read(inputPath + "cyan")
	Y := read(inputPath + "yellow")
	M := read(inputPath + "magenta")
	K := read(inputPath + "black")

	var f []*image.RGBA
	var tempF []*image.NRGBA

	ptX := pt.X
	ptY := pt.Y
	p := image.Point{ptX / outShift, ptY / outShift}
	dst := imaging.New(ptX, ptY, bgColor)
	for x := 0; x < outShift; x++ {
		for y := 0; y < outShift; y++ {
			dstImg := convertAndShiftAllImages(K, C, M, Y, ptX/outShift, 0, 1)
			dst, tempF = combineAllImages(f, dst, dstImg, image.Point{p.X * x, p.Y * y})
			for _, i := range tempF {
				img := setBG(bgColor, toRGBA(i))
				f = append(f, img)
			}
		}
	}
	f = append(f, colorize(f[len(f)-1], bgColor, pt, outShift)...)
	return f
}

func colorize(final *image.RGBA, bgColor color.RGBA, pt image.Point, outShift int) []*image.RGBA {
	var f []*image.RGBA

	ptX := pt.X
	ptY := pt.Y
	var ch []int
	anchors := make(map[int]image.Point)
	for y := 0; y < outShift; y++ {
		for x := 0; x < outShift; x++ {
			ch = append(ch, x+y*outShift)
			anchors[x+y*outShift] = image.Point{(ptX / outShift) * x, (ptY / outShift) * y}
		}
	}

	for i := 0; i < outShift*outShift; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		p := ch[rand.Intn(len(ch))]
		ch = remove(ch, p)
		mask := selectMask(outShift, anchors[p], image.Point{ptX, ptY})
		mask = changeHueToRandom(maskRGBA(final, mask))
		final = combineMaskedWithFinalFrame(final, mask)
		final = setBG(bgColor, final)
		f = append(f, final)
	}
	//f = append(f, final)
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

	C = imaging.Resize(C, s, 0, imaging.Blackman)
	Y = imaging.Resize(Y, s, 0, imaging.Blackman)
	M = imaging.Resize(M, s, 0, imaging.Blackman)
	K = imaging.Resize(K, s, 0, imaging.Blackman)

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

func save(outPath, output string, image image.Image) {
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
}

func random(min, max int) int {
	return min + rand.Intn(max-min)
}
