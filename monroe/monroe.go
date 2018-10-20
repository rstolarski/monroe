package monroe

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Monroe struct {
	rect        image.Rectangle
	inputImage  *image.RGBA
	newImages   [4]*image.CMYK
	alphaMask   *image.Alpha
	outputImage *image.RGBA
}

const minShift = -20
const maxShift = 50

func (m *Monroe) GetOutputImage() *image.RGBA {
	return m.outputImage
}

func NewMonroe(img image.Image) *Monroe {
	m := Monroe{}
	m.inputImage = image.NewRGBA(img.Bounds())
	m.alphaMask = image.NewAlpha(img.Bounds())
	m.outputImage = image.NewRGBA(img.Bounds())
	rect := img.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			m.inputImage.Set(x, y, color.RGBA{
				c.R,
				c.G,
				c.B,
				c.A,
			})
		}
	}

	for i := 0; i < 4; i++ {
		m.newImages[i] = image.NewCMYK(rect)
	}
	m.rect = rect
	return &m
}

func ReadMonroe(input string) *Monroe {
	f, err := os.Open(input + ".png")
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}
	return NewMonroe(img)
}

func (m *Monroe) SplitMonroe() {
	bounds := m.rect
	alpha := color.Alpha{}
	cmyk := color.CMYK{}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//cmyk = color.CMYKModel.Convert(m.inputImage.At(x, y)).(color.CMYK)
			//key := color.CMYKModel.Convert(m.inputImage.At(x, y)).(color.CMYK).K
			cmyk = color.CMYKModel.Convert(m.inputImage.At(x, y)).(color.CMYK)
			alpha = color.AlphaModel.Convert(m.inputImage.At(x, y)).(color.Alpha)

			m.newImages[0].SetCMYK(x, y, color.CMYK{cmyk.C, 0, 0, 0})
			m.newImages[1].SetCMYK(x, y, color.CMYK{0, cmyk.M, 0, 0})
			m.newImages[2].SetCMYK(x, y, color.CMYK{0, 0, cmyk.Y, 0})
			m.newImages[3].SetCMYK(x, y, color.CMYK{0, 0, 0, cmyk.K})
			m.alphaMask.SetAlpha(x, y, color.Alpha{alpha.A})
		}
	}
}

func (m *Monroe) ShiftImages() {
	for i := 0; i < 4; i++ {
		minx, miny, maxx, maxy := generateShifts()
		//maxx, maxy = m.newImages[i].Bounds().Max.X, m.newImages[i].Bounds().Max.Y
		m.newImages[i] = shiftCMYKImage(minx, miny, maxx, maxy, m.newImages[i])
	}
}

func generateShifts() (int, int, int, int) {
	rand.Seed(time.Now().Unix())
	i0 := randInt(minShift, maxShift)
	i1 := randInt(minShift, maxShift)
	i2 := randInt(minShift, maxShift)
	i3 := randInt(minShift, maxShift)
	return i0, i1, i2, i3
}

func shiftCMYKImage(minx, miny, maxx, maxy int, src *image.CMYK) *image.CMYK {
	srcBounds := src.Bounds()
	rect := image.Rect(
		srcBounds.Min.X-minx,
		srcBounds.Min.Y-miny,
		srcBounds.Max.X-maxx,
		srcBounds.Max.Y-maxy,
	)
	dst := image.NewCMYK(rect)
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
	return dst
}

func (m *Monroe) MaskCMYK() {
	bounds := m.rect
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r := color.RGBAModel.Convert(m.outputImage.At(x, y)).(color.RGBA)
			g := color.RGBAModel.Convert(m.outputImage.At(x, y)).(color.RGBA)
			b := color.RGBAModel.Convert(m.outputImage.At(x, y)).(color.RGBA)
			a := color.RGBAModel.Convert(m.inputImage.At(x, y)).(color.RGBA)
			m.outputImage.Set(x, y, color.RGBA{
				r.R,
				g.G,
				b.B,
				a.A})
		}
	}
}

func (m *Monroe) CombineMonroe() {
	bounds := m.outputImage.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c0 := color.CMYKModel.Convert(m.newImages[0].At(x, y)).(color.CMYK)
			c1 := color.CMYKModel.Convert(m.newImages[1].At(x, y)).(color.CMYK)
			c2 := color.CMYKModel.Convert(m.newImages[2].At(x, y)).(color.CMYK)
			c3 := color.CMYKModel.Convert(m.newImages[3].At(x, y)).(color.CMYK)
			//a := color.RGBAModel.Convert(m.inputImage.At(x, y)).(color.RGBA).A
			m.outputImage.Set(x, y, color.CMYK{
				c0.C,
				c1.M,
				c2.Y,
				c3.K})
		}
	}
}

func (m *Monroe) OutputMonroe(image image.Image, output string) {
	f, err := os.Create(output + ".png")
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

func (m *Monroe) OutputMultipleMonroe(output string) {
	for i := 0; i < 4; i++ {
		m.OutputMonroe(m.newImages[i], output+"_"+strconv.Itoa(i))
	}
}

func (m *Monroe) OutputSpecificMonroe(image *image.RGBA, output string) {
	m.OutputMonroe(image, output)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
