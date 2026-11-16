# my_dbms_web

Interface web para o [my_dbms](../my_dbms): um editor SQL no navegador que
conversa com o servidor TCP do banco.

## Como rodar

Em um terminal, o banco como servidor:

```
cd ../my_dbms
go run ./cmd/my_dbms -listen localhost:4321
```

Em outro, a interface:

```
go run .
# abra http://localhost:8080
```

Opções: `-addr` muda o endereço da interface e `-db` o do banco.

## Como funciona

- `main.go`: serve a página e a rota `POST /api/query`. Cada requisição abre
  a sua conexão com o banco e manda os comandos um por linha. Linhas com `--`
  são comentários.
- `client.go`: o cliente do protocolo do my_dbms: uma linha de SQL vai e uma
  linha de JSON (`{"output": ...}` ou `{"error": ...}`) volta.
- `index.html`: editor (Ctrl+Enter executa), lista de tabelas e resultados.
  O SELECT vira uma tabela HTML montada a partir de `columns` e `rows`; os
  outros comandos aparecem como texto.
