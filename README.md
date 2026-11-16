# my_dbms_web

Interface web para o [my_dbms](https://github.com/nataferreiradev/my_dbms):
um editor SQL no navegador que conversa com o servidor TCP do banco.

## Requisitos

Go 1.23 ou mais novo e o my_dbms clonado ao lado desta pasta:

```
git clone https://github.com/nataferreiradev/my_dbms.git
git clone https://github.com/nataferreiradev/my_dbms_web.git
```

A interface usa o servidor TCP do my_dbms (`-listen`).

## Como rodar

```
./dev.sh
# abra http://localhost:8080
```

O script compila e sobe o servidor do my_dbms (`../my_dbms`, porta 4321) e
a interface, e o Ctrl+C derruba os dois. Se a porta do banco já estiver em
uso, por exemplo por um servidor antigo esquecido, ele avisa em vez de a
página falar com a versão velha. `DB_DIR`, `DB_ADDR` e `WEB_ADDR` mudam a
pasta e os endereços.

Para subir cada um à mão:

```
cd ../my_dbms && go run ./cmd/my_dbms -listen localhost:4321
go run .   # em outro terminal
```

## Como funciona

- `main.go`: serve a página e a rota `POST /api/query`. Cada requisição abre
  a sua conexão com o banco. Um comando vai até a linha que termina em `;`,
  então pode ocupar várias linhas, e vai para o banco numa linha só. Linhas
  com `--` são comentários.
- `client.go`: o cliente do protocolo do my_dbms: uma linha de SQL vai e uma
  linha de JSON (`{"output": ...}` ou `{"error": ...}`) volta.
- `index.html`: editor (Ctrl+Enter executa), lista de tabelas e resultados.
  O SELECT vira uma tabela HTML montada a partir de `columns` e `rows`; os
  outros comandos aparecem como texto.

## Estrutura de uma tabela

Clicar numa tabela na lateral mostra a estrutura dela, montada só com
consultas às tabelas do sistema:

- colunas, tipos e `NOT NULL`, de `_catalog`;
- chaves primárias, `UNIQUE` e índices, com páginas e altura da B+tree, de
  `_indexes`;
- linhas, páginas e bytes, de `_tables`;
- o DDL (`CREATE TABLE` e `CREATE INDEX`) que recria a tabela, para copiar.

As tabelas do sistema ficam numa seção separada e mostram os dados direto.
