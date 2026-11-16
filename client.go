package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// timeout limita quanto uma consulta espera pelo servidor.
const timeout = 30 * time.Second

// response é o que o servidor do my_dbms devolve para cada linha.
// Num SELECT, Columns e Rows trazem o resultado; NULL vem como nil.
type response struct {
	Output  string      `json:"output"`
	Error   string      `json:"error"`
	Columns []column    `json:"columns"`
	Rows    [][]*string `json:"rows"`
}

type column struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// conn é uma conexão com o servidor do my_dbms (go run ./cmd/my_dbms -listen ...).
type conn struct {
	c net.Conn
	r *bufio.Reader
}

func dial(addr string) (*conn, error) {
	c, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("conectar ao my_dbms em %s: %w", addr, err)
	}
	return &conn{c: c, r: bufio.NewReader(c)}, nil
}

// send manda um comando e espera a resposta dele.
func (c *conn) send(line string) (response, error) {
	c.c.SetDeadline(time.Now().Add(timeout))
	if _, err := fmt.Fprintln(c.c, line); err != nil {
		return response{}, fmt.Errorf("enviar comando: %w", err)
	}
	data, err := c.r.ReadBytes('\n')
	if err != nil {
		return response{}, fmt.Errorf("ler resposta: %w", err)
	}
	var resp response
	if err := json.Unmarshal(data, &resp); err != nil {
		return response{}, fmt.Errorf("resposta inválida: %w", err)
	}
	return resp, nil
}

func (c *conn) close() {
	c.c.Close()
}
