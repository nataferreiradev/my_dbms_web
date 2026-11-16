package main

import (
	"embed"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"
)

//go:embed index.html
var static embed.FS

type result struct {
	SQL    string `json:"sql"`
	Output string `json:"output"`
	Error  bool   `json:"error"`
}

func main() {
	addr := flag.String("addr", "localhost:8080", "endereço do servidor web")
	dbAddr := flag.String("db", "localhost:4321", "endereço do servidor do my_dbms")
	flag.Parse()

	// O net/http já atende cada requisição numa goroutine, e cada uma abre
	// a sua conexão com o banco: várias abas consultam ao mesmo tempo.
	http.Handle("GET /", http.FileServerFS(static))
	http.HandleFunc("POST /api/query", func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			SQL string `json:"sql"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "json inválido", http.StatusBadRequest)
			return
		}
		writeJSON(w, runAll(*dbAddr, body.SQL))
	})

	log.Printf("usando o my_dbms em %s", *dbAddr)
	log.Printf("abra http://%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// runAll abre uma conexão e manda cada linha não vazia, porque o my_dbms
// entende um comando por linha. Linhas começadas com -- são comentários.
func runAll(addr, sql string) []result {
	results := []result{}
	var c *conn
	for _, line := range strings.Split(sql, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		if line == "exit" || line == "quit" {
			results = append(results, result{line, "aqui não dá para sair: a conexão fecha sozinha", true})
			continue
		}
		if c == nil {
			var err error
			if c, err = dial(addr); err != nil {
				return append(results, result{line, err.Error(), true})
			}
			defer c.close()
		}
		resp, err := c.send(line)
		switch {
		case err != nil:
			return append(results, result{line, err.Error(), true})
		case resp.Error != "":
			results = append(results, result{line, "erro: " + resp.Error, true})
		default:
			results = append(results, result{line, resp.Output, false})
		}
	}
	return results
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
