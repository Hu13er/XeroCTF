package main

import (
	"fmt"
	"image/png"
	"net/http"
	"strings"
)

const (
	flag = "xeroctf{d0nt-oVErth1nk-it!}"
)

type server struct{}

var _ http.Handler = (*server)(nil)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; {
	case strings.HasPrefix(p, "/captcha") && r.Method == "GET":
		s.captcha(w, r)
	case p == "" || p == "/":
		s.index(w, r)
	case strings.HasPrefix(p, "/"):
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404. not found"))
	}
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	t := 0
	c := ""

	if cookie, err := r.Cookie("shampoo"); err == nil {
		v := cookie.Value
		if tt, cc, err := decrypt(v); err == nil {
			t = tt
			c = cc
		}
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("bad request"))
			return
		}

		gotCap := r.Form.Get("val")
		if strings.TrimSpace(gotCap) == strings.TrimSpace(c) {
			t += 1
		}
	}

	f := ""
	if t >= 500 {
		f = "<h2>flag is:" + flag + "</h2>"
		t = 500
	}

	cap := builder.create()
	cookie := encrypt(t, cap.value)
	http.SetCookie(w, &http.Cookie{
		Name:  "shampoo",
		Value: cookie,
	})

	body := `<html>
	<body>
	<h1> CAPTCHA </h1>
	%s
	<p>Try to solve this captcha <b>%d times!</b></p>
	<img src="/captcha?id=%s">
	<form action="/" method="post">
		<label for="val">value</label>
		<input type="text" id="val" name="val"><br>
		<input type="submit" value="Submit">
	</form>
	</body>
</html>
`
	fmt.Fprintf(w, body, f, 500-t, cookie)
}

func (s *server) captcha(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, cap, err := decrypt(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	actualCap, ok := builder.dict[cap]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		return
	}

	w.Header().Add("Content-Type", "image/png")
	png.Encode(w, actualCap.image)
}
