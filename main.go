package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"golang.org/x/net/websocket"

	"github.com/gitxone/gitxone-core/parsers"
)

const defaultHost = ""
const defaultPort = 10098
const bufferSize = 1024 * 8

type Request struct {
	Id      string `json:"id"`
	Command string `json:"command"`
}

type Response struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Err     string `json:"error"`
}

func gitHandlerFactory(path string) func(ws *websocket.Conn) {
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
			tokens = append([]string{"-C", path}, tokens...)

			command := exec.Command("git", tokens...)
			command.Env = []string{
				"PAGER=cat",
			}
			out, err := command.Output()
			res := Response{Id: req.Id, Content: string(out)}
			if err != nil {
				res.Err = err.Error()
			}
			data, _ := json.Marshal(res)
			ws.Write(data)
		}
	}
}

func main() {
	http.HandleFunc(
		"/git",
		func(w http.ResponseWriter, req *http.Request) {
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
			gitHandler := gitHandlerFactory(path)
			s := websocket.Server{Handler: websocket.Handler(gitHandler)}
			s.ServeHTTP(w, req)
		},
	)
	http.Handle(
		"/",
		http.FileServer(http.Dir("dist/")),
	)

	host := flag.String("host", defaultHost, "host address")
	port := flag.Int("port", defaultPort, "port number")
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
