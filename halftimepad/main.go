package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("HLAFTIMEPAD_ADDR")
	if addr == "" {
		addr = "0.0.0.0:8888"
	}
	if err := http.ListenAndServe(addr, &Http{}); err != nil {
		log.Fatalln(err)
	}
}

type Http struct{}

func (*Http) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing remote:", r.RemoteAddr)
	w.Write(HalfTimePath())
}

func init() {
	plain = []byte(`Hola dear friends.
this is a sample p14in t3xt and if u are readin this, it means u already broke it.
congrats.
here is ur flag  xeroctf{d0nt-use_1timepad_haf1y.he!shampoo_-_-_}
...and lets try to continue this text a little bit more.
ok i think its enough!
_*Good Luck*_
^___^`)
}
