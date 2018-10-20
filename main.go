package main

import (
	"github.com/rtropisz/monroe/monroe"
)

func main() {
	m := monroe.ReadMonroe("input")
	m.SplitMonroe()
	m.ShiftImages()

	m.CombineMonroe()
	m.MaskCMYK()
	//m.OutputMultipleMonroe("output")
	m.OutputMonroe(m.GetOutputImage(), "output")
}
