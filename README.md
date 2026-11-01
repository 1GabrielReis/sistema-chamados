# Sistema de Chamados

Aplicação web para abertura, distribuição e acompanhamento de chamados internos. O projeto foi desenvolvido em Go, React e PostgreSQL e pode ser executado integralmente com Docker Compose.

## Funcionalidades

- cadastro, visualização, edição, busca e filtragem de chamados;
- prioridades baixa, média e alta;
- status aberto, em andamento, resolvido e fechado;
- atribuição manual de um responsável;
- distribuição automática para o responsável com menos chamados não concluídos;
- cadastro e consulta de pessoas por código público;
- atualização do nome de uma pessoa;
- habilitação de uma pessoa como responsável mediante autorização;
- interface responsiva para computador e celular;
- três responsáveis iniciais disponíveis no banco.

## Tecnologias

- Go 1.27 e biblioteca `net/http` para a API;
- PostgreSQL 18 para persistência;
- `pgx`/`pgxpool` para acesso ao banco;
- React 19 e Vite 8 para o frontend;
- Lucide React para os ícones;
- Nginx para servir o frontend e encaminhar `/api` para a API;
- Docker e Docker Compose para execução do ambiente completo.

## Arquitetura

O backend separa as responsabilidades em camadas:

```text
Requisição HTTP
      ↓
Handler → Service → Repository → PostgreSQL
      ↑
Resposta JSON
```

- `cmd/server`: composição e inicialização da API;
- `internal/domain`: entidades do domínio;
- `internal/handler`: entrada e saída HTTP;
- `internal/service`: validações e regras de negócio;
- `internal/repository`: contratos de persistência;
- `internal/repository/postgres`: implementação PostgreSQL;
- `internal/database`: criação do pool de conexões;
- `migrations`: criação das tabelas, índices e dados iniciais;
- `frontend`: aplicação React.

O frontend conversa somente com o caminho relativo `/api`. Em desenvolvimento, o Vite encaminha esse caminho para `localhost:8080`. No Docker, o Nginx encaminha para o serviço `api` pela rede interna do Compose.

Um único `Dockerfile` utiliza build multi-stage. O Compose seleciona o estágio final `api` para o backend e o estágio final `frontend` para a interface, produzindo duas imagens independentes a partir do mesmo arquivo.

## Executar com Docker

### Pré-requisitos

- Docker Desktop com Docker Compose.

### 1. Configure o ambiente

No PowerShell:

```powershell
Copy-Item .env.example .env
```

Abra `.env` e substitua `troque_esta_senha` e `SUA_SENHA` pela mesma senha. O arquivo `.env` não deve ser enviado ao GitHub.

### 2. Suba a aplicação

```powershell
docker compose up --build -d
```

Na primeira inicialização de um volume vazio, o PostgreSQL executa automaticamente as migrations de criação e inclui os responsáveis `RESP-001`, `RESP-002` e `RESP-003`.

### 3. Acesse

- aplicação: http://localhost:3000
- API: http://localhost:8080
- verificação da API: http://localhost:8080/health

Para acompanhar os contêineres:

```powershell
docker compose ps
docker compose logs -f
```

Para parar sem apagar os dados:

```powershell
docker compose down
```

O comando abaixo também apaga o volume e todos os dados locais do banco. Use apenas quando quiser recomeçar:

```powershell
docker compose down -v
```

## Desenvolvimento sem o frontend em Docker

Com o banco em execução, defina a conexão e inicie a API:

```powershell
$env:DATABASE_URL='postgres://app:SUA_SENHA@localhost:5432/sistema_chamados?sslmode=disable'
go run ./cmd/server
```

Em outro terminal, inicie o frontend:

```powershell
Set-Location frontend
pnpm install
pnpm dev
```

O Vite disponibiliza a interface em http://localhost:5173.

## Testes

Backend:

```powershell
go test ./...
```

Os testes de integração dos repositórios utilizam `DATABASE_URL`. Quando essa variável não está definida, eles são ignorados.

Frontend:

```powershell
Set-Location frontend
pnpm test
pnpm build
```

## Decisões do domínio

- `aberto` e `em_andamento` são considerados não concluídos na distribuição automática;
- `resolvido` e `fechado` são considerados concluídos;
- a pessoa utiliza um código público para abrir chamados, sem autenticação nesta primeira versão;
- um responsável já habilitado pode habilitar outra pessoa;
- a distribuição automática é decidida pelo backend, evitando regras divergentes entre clientes;
- as chaves estrangeiras usam `ON DELETE RESTRICT` para impedir chamados órfãos.

## Referências e bibliotecas externas

- [Go](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- [pgx](https://github.com/jackc/pgx)
- [React](https://react.dev/)
- [Vite](https://vite.dev/)
- [Lucide](https://lucide.dev/)
- [Nginx](https://nginx.org/)
- [Docker Compose](https://docs.docker.com/compose/)
