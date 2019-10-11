package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/net/websocket"

	"github.com/gitxone/gitxone-core/completion"
	"github.com/gitxone/gitxone-core/parsers"

	"github.com/rakyll/statik/fs"

	_ "github.com/gitxone/gitxone-core/statik"
)

var gitCommand = "git"
var host = ""
var port = 10098

const bufferSize = 1024 * 8

type Request struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

type Response struct {
	Type    string      `json:"type"`
	Tokens  []string    `json:"tokens"`
	Content interface{} `json:"content"`
	Err     string      `json:"error"`
}

func gitHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	paths, ok := params["path"]
	if !ok || len(paths) != 1 {
		http.NotFound(w, req)
		return
	}
	path := paths[0]
	if path[0] == '~' {
		path = fmt.Sprintf("%s%s", os.Getenv("HOME"), path[1:])
	} else {
		path = fmt.Sprintf("/%s", path)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.NotFound(w, req)
		return
	}
	gitHandler := gitSocketHandlerFactory(path)
	s := websocket.Server{Handler: websocket.Handler(gitHandler)}
	s.ServeHTTP(w, req)
}

func gitSocketHandlerFactory(path string) func(ws *websocket.Conn) {
	// just injecting path into websocket handler.
	return func(ws *websocket.Conn) {
		args := make([]byte, bufferSize)
		for {
			n, err := ws.Read(args)
			if err != nil {
				return
			}
			var req Request
			json.Unmarshal(args[:n], &req)

			tokens := parsers.ParseCommand(req.Command)
			res := Response{Type: req.Type, Tokens: tokens}
			switch req.Type {
			case "exec":
				realTokens := append([]string{"-C", path}, tokens...)

				cmd := exec.Command(gitCommand, realTokens...)
				cmd.Env = []string{
					"PAGER=cat",
				}
				out, err := cmd.CombinedOutput()

				res.Content = string(out)
				if err != nil {
					res.Err = err.Error()
				}
			case "complete":
				cands := completion.Complete(path, tokens)
				res.Content = cands
			}
			data, _ := json.Marshal(res)
			ws.Write(data)
		}
	}
}

func main() {
	http.HandleFunc("/git", gitHandler)

	statikFS, _ := fs.New()
	http.Handle("/_nuxt/", http.FileServer(statikFS))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := statikFS.Open("/index.html")
		if err != nil {
			log.Fatal(err)
		}
		fi, _ := f.Stat()
		data := make([]byte, fi.Size())
		f.Read(data)
		w.Write(data)
	})

	host := *flag.String("host", host, "host address")
	port := *flag.Int("port", port, "port number")
	gitCommand = *flag.String("git", gitCommand, "Git path")

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
