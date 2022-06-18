package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	flag = `xeroctf{"Lumos" is a magical spell from Harry Potter films and books}`
)

type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; {
	case p == "" && r.Method == "GET":
		fallthrough
	case p == "/" && r.Method == "GET":
		fallthrough
	case strings.HasPrefix(p, "/help") && r.Method == "GET":
		s.help(w, r)
	case strings.HasPrefix(p, "/login") && r.Method == "POST":
		s.login(w, r)
	case strings.HasPrefix(p, "/desc") && r.Method == "GET":
		s.desc(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found: see /help"))
	}
}

type loginRequest struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type loginResponse struct {
	Flag string `json:"flag"`
}

func (*Server) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := login(req.Name, req.Pass)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(loginResponse{
			Flag: "https://www.youtube.com/shorts/5rzdwq78qWc",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{
		Flag: flag,
	})
}

type descResponse struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (*Server) desc(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	log.Println("Desc", name)

	nds, err := desc(name)
	if err != nil {
		log.Println("ERR", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var result []descResponse
	for _, nd := range nds {
		result = append(result, descResponse{
			Name: nd.name,
			Desc: nd.desc,
		})
	}

	json.NewEncoder(w).Encode(result)
}

type helpResponse struct {
	Msg string `json:"msg"`
}

func (*Server) help(w http.ResponseWriter, r *http.Request) {
	msg := `Available Routes:
	*) GET /help: This help
	*) GET /desc?name=<NAME>
		output:
			json{
				name ~> <Name of User>
				desc ~> <Desc of User>
			}
	*) POST /login:
		input:
			json{
				name ~> <Name of User>
				pass ~> <Pass of User>
			}
		output:
			json{
				flag ~> <Ur most wanted flag>
			}
	`
	json.NewEncoder(w).Encode(helpResponse{Msg: msg})
}

func main() {
	s := &Server{}
	addr := os.Getenv("POTTER_ADDR")
	if addr == "" {
		addr = "0.0.0.0:7887"
	}
	if err := http.ListenAndServe(addr, s); err != nil {
		log.Fatalln(err)
	}
}
