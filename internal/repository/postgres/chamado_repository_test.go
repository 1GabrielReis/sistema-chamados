package postgres

import (
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/1GabrielReis/sistema-chamados/internal/database"
	"github.com/1GabrielReis/sistema-chamados/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

func abrirBancoTeste(t *testing.T) (context.Context, *pgxpool.Pool) {
	t.Helper()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL não definida")
	}

	ctx := context.Background()

	pool, err := database.Conectar(ctx, databaseURL)
	if err != nil {
		t.Fatalf("não foi possível conectar ao banco: %v", err)
	}
	t.Cleanup(pool.Close)

	return ctx, pool
}

func novoPessoaTeste() domain.Pessoa {
	return domain.Pessoa{
		RegistroPublico: "TEST-" + rand.Text(),
		Nome:            "Pessoa de teste",
	}
}

func novoChamadoTeste(solicitante domain.Pessoa, responsavel domain.Pessoa) domain.Chamado {
	return domain.Chamado{
		Solicitante: solicitante,
		Responsavel: responsavel,
		Titulo:      "TEST-" + rand.Text(),
		Descricao:   "chamado de teste",
		Prioridade:  domain.PrioridadeMedia,
		Status:      domain.StatusAberto,
	}
}

func criarSolicitante(t *testing.T, ctx context.Context, repositorio *PessoaRepository) domain.Pessoa {
	t.Helper()

	solicitante := novoPessoaTeste()
	solicitante.Nome = "Solicitante de teste"

	solicitanteCriando, err := repositorio.Criar(ctx, solicitante)
	if err != nil {
		t.Fatalf("não foi possível criar solicitante: %v", err)
	}
	return solicitanteCriando
}

func criarResponsavel(t *testing.T, ctx context.Context, repositorio *PessoaRepository) domain.Pessoa {
	t.Helper()

	responsavelCriando := criarSolicitante(t, ctx, repositorio)
	responsavelCriando.Nome = "Responsavel de teste"
	responsavelCriando.EhResponsavel = true

	responsavelAtualizado, err := repositorio.Atualizar(ctx, responsavelCriando)
	if err != nil {
		t.Fatalf("não foi possível criar Resposanvel: %v", err)
	}

	return responsavelAtualizado
}

func criarChamado(t *testing.T, ctx context.Context, repositorio *ChamadoRepository,
	solicitante domain.Pessoa, responsavel domain.Pessoa) domain.Chamado {
	t.Helper()

	chamado := novoChamadoTeste(solicitante, responsavel)

	chamadoCriado, err := repositorio.Criar(ctx, chamado)
	if err != nil {
		t.Fatalf("não foi possível criar chamado: %v", err)
	}

	return chamadoCriado
}

func registrarLimpezaTeste(
	t *testing.T,
	pool *pgxpool.Pool,
	pessoaIDs *[]int64,
	chamadoIDs *[]int64,
) {

	t.Helper()
	t.Cleanup(func() {
		if len(*chamadoIDs) > 0 {
			_, err := pool.Exec(
				context.Background(),
				"DELETE FROM chamados WHERE id = ANY($1::bigint[])",
				*chamadoIDs,
			)
			if err != nil {
				t.Errorf("não foi possível remover chamados de teste: %v", err)
			}
		}

		if len(*pessoaIDs) > 0 {
			_, err := pool.Exec(
				context.Background(),
				"DELETE FROM pessoas WHERE id = ANY($1::bigint[])",
				*pessoaIDs,
			)
			if err != nil {
				t.Errorf("não foi possível remover pessoas de teste: %v", err)
			}
		}
	})
}

func TestChamadoRepository_Criar(t *testing.T) {
	ctx, pool := abrirBancoTeste(t)

	pessoaRepository := NovoPessoaRepository(pool)
	chamadoRepository := NovoChamadoRepository(pool)

	pessoaIDs := make([]int64, 0)
	chamadoIDs := make([]int64, 0)

	registrarLimpezaTeste(t, pool, &pessoaIDs, &chamadoIDs)

	solicitante := criarSolicitante(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, solicitante.ID)

	responsavel := criarResponsavel(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, responsavel.ID)

	chamado := novoChamadoTeste(solicitante, responsavel)

	chamadoCriado, err := chamadoRepository.Criar(ctx, chamado)
	if err != nil {
		t.Fatalf("não foi possível criar chamado: %v", err)
	}
	chamadoIDs = append(chamadoIDs, chamadoCriado.ID)

	if chamadoCriado.ID <= 0 {
		t.Errorf("esperava ID positivo, recebeu %d", chamadoCriado.ID)
	}

	if chamadoCriado.Solicitante.ID != solicitante.ID {
		t.Errorf(
			"solicitante diferente: esperado ID %d, recebido ID %d",
			solicitante.ID,
			chamadoCriado.Solicitante.ID,
		)
	}

	if chamadoCriado.Solicitante.EhResponsavel {
		t.Error("esperava solicitante como pessoa comum")
	}

	if chamadoCriado.Responsavel.ID != responsavel.ID {
		t.Errorf(
			"responsável diferente: esperado ID %d, recebido ID %d",
			responsavel.ID,
			chamadoCriado.Responsavel.ID,
		)
	}

	if !chamadoCriado.Responsavel.EhResponsavel {
		t.Error("esperava pessoa responsável")
	}

	if chamadoCriado.Titulo != chamado.Titulo {
		t.Errorf("titulo diferente: esperado %q, recebido %q",
			chamado.Titulo, chamadoCriado.Titulo)
	}

	if chamadoCriado.Descricao != chamado.Descricao {
		t.Errorf("discricao diferente: esperado %q, recebido %q",
			chamado.Descricao, chamadoCriado.Descricao)
	}

	if chamadoCriado.Prioridade != chamado.Prioridade {
		t.Errorf("prioridade diferente: esperado %q, recebido %q",
			chamado.Prioridade, chamadoCriado.Prioridade)
	}

	if chamadoCriado.Status != chamado.Status {
		t.Errorf("status diferente: esperado %q, recebido %q",
			chamado.Status, chamadoCriado.Status)
	}

	if chamadoCriado.DataAbertura.IsZero() {
		t.Error("esperava data de abertura preenchida")
	}

}

func TestChamadoRepository_BuscarPorID(t *testing.T) {
	ctx, pool := abrirBancoTeste(t)

	pessoaRepository := NovoPessoaRepository(pool)
	chamadoRepository := NovoChamadoRepository(pool)

	pessoaIDs := make([]int64, 0)
	chamadoIDs := make([]int64, 0)

	registrarLimpezaTeste(t, pool, &pessoaIDs, &chamadoIDs)

	solicitante := criarSolicitante(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, solicitante.ID)

	responsavel := criarResponsavel(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, responsavel.ID)

	chamado := criarChamado(t, ctx, chamadoRepository, solicitante, responsavel)
	chamadoIDs = append(chamadoIDs, chamado.ID)

	chamadoBuscado, err := chamadoRepository.BuscarPorID(ctx, chamado.ID)
	if err != nil {
		t.Fatalf("não foi possível buscar chamado: %v", err)
	}

	if chamadoBuscado.ID != chamado.ID {
		t.Errorf(
			"ID diferente: esperado %d, recebido %d",
			chamado.ID,
			chamadoBuscado.ID,
		)
	}

	if chamadoBuscado.Solicitante.ID != solicitante.ID {
		t.Errorf(
			"solicitante diferente: esperado ID %d, recebido ID %d",
			solicitante.ID,
			chamadoBuscado.Solicitante.ID,
		)
	}

	if chamadoBuscado.Solicitante.EhResponsavel {
		t.Error("esperava solicitante como pessoa comum")
	}

	if chamadoBuscado.Responsavel.ID != responsavel.ID {
		t.Errorf(
			"responsável diferente: esperado ID %d, recebido ID %d",
			responsavel.ID,
			chamadoBuscado.Responsavel.ID,
		)
	}

	if !chamadoBuscado.Responsavel.EhResponsavel {
		t.Error("esperava pessoa responsável")
	}

	if chamadoBuscado.Titulo != chamado.Titulo {
		t.Errorf("titulo diferente: esperado %q, recebido %q",
			chamado.Titulo, chamadoBuscado.Titulo)
	}

	if chamadoBuscado.Descricao != chamado.Descricao {
		t.Errorf("discricao diferente: esperado %q, recebido %q",
			chamado.Descricao, chamadoBuscado.Descricao)
	}

	if chamadoBuscado.Prioridade != chamado.Prioridade {
		t.Errorf("prioridade diferente: esperado %q, recebido %q",
			chamado.Prioridade, chamadoBuscado.Prioridade)
	}

	if chamadoBuscado.Status != chamado.Status {
		t.Errorf("status diferente: esperado %q, recebido %q",
			chamado.Status, chamadoBuscado.Status)
	}

	if chamadoBuscado.DataAbertura.IsZero() {
		t.Error("esperava data de abertura preenchida")
	}

}

func TestChamadoRepository_Listar(t *testing.T) {
	ctx, pool := abrirBancoTeste(t)

	pessoaRepository := NovoPessoaRepository(pool)
	chamadoRepository := NovoChamadoRepository(pool)

	pessoaIDs := make([]int64, 0)
	chamadoIDs := make([]int64, 0)

	registrarLimpezaTeste(t, pool, &pessoaIDs, &chamadoIDs)

	solicitante1 := criarSolicitante(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, solicitante1.ID)

	solicitante2 := criarSolicitante(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, solicitante2.ID)

	responsavel1 := criarResponsavel(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, responsavel1.ID)

	responsavel2 := criarResponsavel(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, responsavel2.ID)

	chamado1 := criarChamado(t, ctx, chamadoRepository, solicitante1, responsavel1)
	chamadoIDs = append(chamadoIDs, chamado1.ID)

	chamado2 := criarChamado(t, ctx, chamadoRepository, solicitante1, responsavel2)
	chamadoIDs = append(chamadoIDs, chamado2.ID)

	chamado3 := criarChamado(t, ctx, chamadoRepository, solicitante2, responsavel1)
	chamadoIDs = append(chamadoIDs, chamado3.ID)

	chamado4 := criarChamado(t, ctx, chamadoRepository, solicitante2, responsavel2)
	chamadoIDs = append(chamadoIDs, chamado4.ID)

	esperados := map[int64]bool{
		chamado1.ID: false,
		chamado2.ID: false,
		chamado3.ID: false,
		chamado4.ID: false,
	}

	chamados, err := chamadoRepository.Listar(ctx)
	if err != nil {
		t.Fatalf("não foi possível listar chamados: %v", err)
	}

	for _, chamado := range chamados {

		_, pertenceAoTeste := esperados[chamado.ID]
		if !pertenceAoTeste {
			continue
		}

		esperados[chamado.ID] = true

		if chamado.Solicitante.EhResponsavel {
			t.Error("esperava solicitante como pessoa comum")
		}

		if !chamado.Responsavel.EhResponsavel {
			t.Error("esperava pessoa responsável")
		}

		if chamado.DataAbertura.IsZero() {
			t.Error("esperava data de abertura preenchida")
		}

	}

	for id, encontrado := range esperados {
		if !encontrado {
			t.Errorf("chamado de teste não apareceu na listagem: id=%d", id)
		}
	}

}

func TestChamadoRepository_Atualizar(t *testing.T) {
	ctx, pool := abrirBancoTeste(t)

	pessoaRepository := NovoPessoaRepository(pool)
	chamadoRepository := NovoChamadoRepository(pool)

	pessoaIDs := make([]int64, 0)
	chamadoIDs := make([]int64, 0)

	registrarLimpezaTeste(t, pool, &pessoaIDs, &chamadoIDs)

	solicitante := criarSolicitante(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, solicitante.ID)

	responsavel1 := criarResponsavel(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, responsavel1.ID)

	responsavel2 := criarResponsavel(t, ctx, pessoaRepository)
	pessoaIDs = append(pessoaIDs, responsavel2.ID)

	chamado := criarChamado(t, ctx, chamadoRepository, solicitante, responsavel1)
	chamadoIDs = append(chamadoIDs, chamado.ID)

	dataAberturaOriginal := chamado.DataAbertura
	chamado.Responsavel = responsavel2
	chamado.Titulo = "Atuaalização -" + rand.Text()
	chamado.Descricao = "atualização de chamado"
	chamado.Prioridade = domain.PrioridadeBaixa
	chamado.Status = domain.StatusEmAndamento

	chamadoAtualizado, err := chamadoRepository.Atualizar(ctx, chamado)
	if err != nil {
		t.Fatalf("não foi possível atualizar chamado: %v", err)
	}

	chamadoPersistido, err := chamadoRepository.BuscarPorID(ctx, chamadoAtualizado.ID)
	if err != nil {
		t.Fatalf("não foi possível buscar chamado atualizado: %v", err)
	}

	if chamadoAtualizado.ID != chamado.ID {
		t.Errorf("ID diferente: esperado %d, recebido %d",
			chamado.ID,
			chamadoAtualizado.ID)
	}

	if chamadoPersistido.Solicitante.ID != solicitante.ID {
		t.Errorf(
			"solicitante diferente: esperado ID %d, recebido ID %d",
			chamado.Solicitante.ID,
			chamadoPersistido.Solicitante.ID,
		)
	}

	if chamadoPersistido.Solicitante.EhResponsavel {
		t.Error("esperava solicitante como pessoa comum")
	}

	if chamadoPersistido.Responsavel.ID != responsavel2.ID {
		t.Errorf(
			"responsável diferente: esperado ID %d, recebido ID %d",
			responsavel2.ID,
			chamadoPersistido.Responsavel.ID,
		)
	}

	if !chamadoPersistido.Responsavel.EhResponsavel {
		t.Error("esperava pessoa responsável")
	}

	if chamado.Titulo != chamadoPersistido.Titulo {
		t.Errorf("titulo diferente: esperado %q, recebido %q",
			chamado.Titulo, chamadoPersistido.Titulo)
	}

	if chamado.Descricao != chamadoPersistido.Descricao {
		t.Errorf("discricao diferente: esperado %q, recebido %q",
			chamado.Descricao, chamadoPersistido.Descricao)
	}

	if chamado.Prioridade != chamadoPersistido.Prioridade {
		t.Errorf("prioridade diferente: esperado %q, recebido %q",
			chamado.Prioridade, chamadoPersistido.Prioridade)
	}

	if chamado.Status != chamadoPersistido.Status {
		t.Errorf("status diferente: esperado %q, recebido %q",
			chamado.Status, chamadoPersistido.Status)
	}

	if !chamadoPersistido.DataAbertura.Equal(dataAberturaOriginal) {
		t.Errorf(
			"data de abertura foi alterada: esperada %v, recebida %v",
			dataAberturaOriginal,
			chamadoPersistido.DataAbertura,
		)
	}
}
