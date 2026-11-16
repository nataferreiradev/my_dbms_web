package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	dbDir := flag.String("db", "../my_dbms", "pasta do projeto my_dbms")
	flag.Parse()

	dir, err := filepath.Abs(*dbDir)
	if err != nil {
		log.Fatal(err)
	}

	bin, err := build(dir)
	if err != nil {
		log.Fatal(err)
	}

	r := &repl{bin: bin, dir: dir}
	if err := r.start(); err != nil {
		log.Fatal(err)
	}
	defer r.stop()

	http.Handle("GET /", http.FileServerFS(static))
	http.HandleFunc("POST /api/query", func(w http.ResponseWriter, req *http.Request) {
		var body struct {
			SQL string `json:"sql"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			http.Error(w, "json inválido", http.StatusBadRequest)
			return
		}
		writeJSON(w, runAll(r, body.SQL))
	})

	log.Printf("my_dbms em %s", dir)
	log.Printf("abra http://%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// build compila o REPL do my_dbms para um executável temporário.
func build(dir string) (string, error) {
	bin := filepath.Join(os.TempDir(), "my_dbms_web_repl")
	if filepath.Separator == '\\' {
		bin += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/my_dbms")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("compilar my_dbms em %s: %w\n%s", dir, err, out)
	}
	return bin, nil
}

// runAll manda cada linha não vazia para o REPL, que entende um comando
// por linha. Linhas começadas com -- são comentários.
func runAll(r *repl, sql string) []result {
	results := []result{}
	for _, line := range strings.Split(sql, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		if line == "exit" || line == "quit" {
			results = append(results, result{line, "aqui não dá para sair do REPL", true})
			continue
		}
		out, err := r.run(line)
		if err != nil {
			results = append(results, result{line, err.Error(), true})
			continue
		}
		results = append(results, result{line, out, strings.HasPrefix(out, "erro:")})
	}
	return results
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
