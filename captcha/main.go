package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("CAPTCHA_ADDR")
	if addr == "" {
		addr = "0.0.0.0:8000"
	}
	dir := os.Getenv("CAPTCHA_DIR")
	if dir == "" {
		dir = "captcha/"
	}
	builder.load(dir)

	if err := http.ListenAndServe(addr, &server{}); err != nil {
		log.Fatalln(err)
	}
}
