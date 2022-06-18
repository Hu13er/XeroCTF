package main

import (
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"

	"golang.org/x/image/bmp"
)

var (
	imgfilename    = ""
	innerfilename  = ""
	outputfilename = ""
	encode         = true
)

func readArgs() {
	usage := func() {
		log.Fatalln("Usage: \n" +
			"\twm encode <IMG_FILE> <INNER_IMG> <OUTPUT>\n" +
			"\twm decode <IMG_FILE> <OUTPUT>")
	}
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "encode":
		encode = true
	case "decode":
		encode = false
	default:
		usage()
	}
	imgfilename = os.Args[2]
	if encode {
		if len(os.Args) < 5 {
			usage()
		}
		innerfilename = os.Args[3]
		outputfilename = os.Args[4]
	} else {
		if len(os.Args) < 4 {
			usage()
		}
		outputfilename = os.Args[3]
	}
}

func main() {
	readArgs()

	imgf, err := os.Open(imgfilename)
	if err != nil {
		log.Fatalf("Could not open file %q: %v\n", imgfilename, err)
	}
	defer imgf.Close()

	img, _, err := image.Decode(imgf)
	if err != nil {
		log.Fatalln("Could not decode image:", err)
	}

	outputf, err := os.Create(outputfilename)
	if err != nil {
		log.Fatalf("Could not write file %q: %v\n", outputfilename, err)
	}
	defer outputf.Close()

	if encode {
		innerf, err := os.Open(innerfilename)
		if err != nil {
			log.Fatalf("Could not open file %q: %v\n", innerfilename, err)
		}
		defer innerf.Close()

		inner, _, err := image.Decode(innerf)
		if err != nil {
			log.Fatalln("Could not decode inner image:", err)
		}

		output, err := watermark(img, inner)
		if err != nil {
			log.Fatalln("Could not watermark:", err)
		}

		err = bmp.Encode(outputf, output)
		if err != nil {
			log.Fatalf("Could not write file %q: %v\n", outputfilename, err)
		}
	} else {
		output, err := dewatermark(img)
		if err != nil {
			log.Fatalln("Could not dewatermark:", err)
		}

		err = bmp.Encode(outputf, output)
		if err != nil {
			log.Fatalf("Could not write file %q: %v\n", outputfilename, err)
		}
	}
}

func watermark(img image.Image, inner image.Image) (image.Image, error) {
	imgBound := img.Bounds()
	innerBound := img.Bounds()

	if innerBound.Max.X > imgBound.Max.X ||
		innerBound.Max.Y > imgBound.Max.Y {
		return nil, errors.New("orginal image size should be bigger than inner image")
	}

	output := image.NewGray(imgBound)
	for x := 0; x < imgBound.Max.X; x++ {
		for y := 0; y < imgBound.Max.Y; y++ {
			g := makeGray(img.At(x, y))
			g.Y >>= 1
			g.Y <<= 1
			if x < innerBound.Max.X &&
				y < innerBound.Max.Y {
				g2 := makeGray(inner.At(x, y))
				if g2.Y >= 256/2 {
					g.Y |= 1
				}
			}
			output.SetGray(x, y, g)
		}
	}

	return output, nil
}

func dewatermark(img image.Image) (image.Image, error) {
	output := image.NewGray(img.Bounds())
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.X; y++ {
			r, _, _, _ := img.At(x, y).RGBA()
			var g uint8
			if r&1 > 0 {
				g = 255
			} else {
				g = 0
			}
			output.SetGray(x, y, color.Gray{g})
		}
	}
	return output, nil
}

func makeGray(c color.Color) color.Gray {
	rr, gg, bb, _ := c.RGBA()
	r := math.Pow(float64(rr), 2.2)
	g := math.Pow(float64(gg), 2.2)
	b := math.Pow(float64(bb), 2.2)
	m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
	Y := uint16(m + 0.5)
	return color.Gray{uint8(Y >> 8)}
}

func init() {
	log.SetFlags(0)
}
