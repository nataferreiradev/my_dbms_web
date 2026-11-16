package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

// prompt é o que o REPL do my_dbms imprime quando está pronto para o
// próximo comando. É ele que marca o fim da saída de cada comando.
const prompt = "my_dbms> "

// repl mantém um processo do my_dbms aberto e manda um comando por vez.
type repl struct {
	bin string // executável do my_dbms
	dir string // pasta onde ele roda (onde ficam my_dbms.conf e data/)

	mu  sync.Mutex
	cmd *exec.Cmd
	in  io.WriteCloser
	out *bufio.Reader
}

func (r *repl) start() error {
	cmd := exec.Command(r.bin)
	cmd.Dir = r.dir
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("iniciar my_dbms: %w", err)
	}
	r.cmd, r.in, r.out = cmd, in, bufio.NewReader(out)

	// Descarta o primeiro prompt.
	if _, err := r.readUntilPrompt(); err != nil {
		r.stop()
		return fmt.Errorf("my_dbms não respondeu: %w", err)
	}
	return nil
}

func (r *repl) stop() {
	if r.cmd == nil {
		return
	}
	r.in.Close()
	r.cmd.Process.Kill()
	r.cmd.Wait()
	r.cmd = nil
}

// readUntilPrompt lê a saída até o REPL imprimir o prompt de novo.
func (r *repl) readUntilPrompt() (string, error) {
	var buf bytes.Buffer
	for !bytes.HasSuffix(buf.Bytes(), []byte(prompt)) {
		b, err := r.out.ReadByte()
		if err != nil {
			return buf.String(), err
		}
		buf.WriteByte(b)
	}
	return strings.TrimRight(strings.TrimSuffix(buf.String(), prompt), "\n"), nil
}

// run manda uma linha para o REPL e devolve o que ele respondeu.
// Se o processo morreu, sobe outro na próxima chamada.
func (r *repl) run(line string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.cmd == nil {
		if err := r.start(); err != nil {
			return "", err
		}
	}
	if _, err := io.WriteString(r.in, line+"\n"); err != nil {
		r.stop()
		return "", fmt.Errorf("enviar comando: %w", err)
	}
	out, err := r.readUntilPrompt()
	if err != nil {
		r.stop()
		return out, fmt.Errorf("ler resposta: %w", err)
	}
	return out, nil
}
