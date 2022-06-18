package main

import (
	"image"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
)

var (
	builder = &captchaFactory{
		set:  make([]*captcha, 0),
		dict: map[string]*captcha{},
	}
)

type captcha struct {
	value string
	image image.Image
}

type captchaFactory struct {
	set  []*captcha
	dict map[string]*captcha
}

func (cf *captchaFactory) load(dir string) error {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		// ignore not png files
		if !strings.HasSuffix(f.Name(), ".png") {
			continue
		}
		p := path.Join(dir, f.Name())
		b, err := os.Open(p)
		if err != nil {
			return err
		}
		img, err := png.Decode(b)
		if err != nil {
			return err
		}

		cap := &captcha{
			value: strings.TrimSuffix(f.Name(), ".png"),
			image: img,
		}
		cf.set = append(cf.set, cap)
		cf.dict[cap.value] = cap
	}
	return nil
}

func (cf *captchaFactory) create() *captcha {
	n := rand.Intn(len(cf.set))
	return cf.set[n]
}
