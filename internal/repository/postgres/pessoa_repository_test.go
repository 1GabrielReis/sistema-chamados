package postgres

import (
	"context"
	"crypto/rand"
	"os"
	"testing"

	"github.com/1GabrielReis/sistema-chamados/internal/database"
	"github.com/1GabrielReis/sistema-chamados/internal/domain"
)

func TestPessoaRepository_Criar(t *testing.T) {
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

	repositorio := NovoPessoaRepository(pool)

	pessoa := domain.Pessoa{
		RegistroPublico: "TEST-" + rand.Text(),
		Nome:            "Pessoa de teste",
	}

	pessoaCriada, err := repositorio.Criar(ctx, pessoa)
	if err != nil {
		t.Fatalf("não foi possível criar pessoa: %v", err)
	}

	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM pessoas WHERE id = $1",
			pessoaCriada.ID,
		)
		if err != nil {
			t.Errorf("não foi possível remover pessoa de teste: %v", err)
		}
	})

	if pessoaCriada.ID <= 0 {
		t.Errorf("esperava ID positivo, recebeu %d", pessoaCriada.ID)
	}

	if pessoaCriada.RegistroPublico != pessoa.RegistroPublico {
		t.Errorf("registro público diferente: esperado %q, recebido %q",
			pessoa.RegistroPublico, pessoaCriada.RegistroPublico)
	}

	if pessoaCriada.Nome != pessoa.Nome {
		t.Errorf("nome diferente: esperado %q, recebido %q",
			pessoa.Nome, pessoaCriada.Nome)
	}

	if pessoaCriada.EhResponsavel {
		t.Error("esperava pessoa não responsável")
	}

	if pessoaCriada.DataCadastro.IsZero() {
		t.Error("esperava data de cadastro preenchida")
	}

}

func TestPessoaRepository_BuscarPorCodigoPublico(t *testing.T) {
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

	repositorio := NovoPessoaRepository(pool)

	pessoa := domain.Pessoa{
		RegistroPublico: "TEST-" + rand.Text(),
		Nome:            "Pessoa de teste",
	}

	pessoaCriada, err := repositorio.Criar(ctx, pessoa)
	if err != nil {
		t.Fatalf("não foi possível criar pessoa: %v", err)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM pessoas WHERE id = $1",
			pessoaCriada.ID,
		)
		if err != nil {
			t.Errorf("não foi possível remover pessoa de teste: %v", err)
		}
	})

	buscarPessoa, err := repositorio.BuscarPorCodigoPublico(ctx, pessoa.RegistroPublico)
	if err != nil {
		t.Fatalf("não foi possível buscar pessoa: %v", err)
	}

	if buscarPessoa.RegistroPublico != pessoaCriada.RegistroPublico {
		t.Errorf(
			"registro público diferente: esperado %q, recebido %q",
			pessoaCriada.RegistroPublico,
			buscarPessoa.RegistroPublico,
		)
	}

	if buscarPessoa.Nome != pessoaCriada.Nome {
		t.Errorf("nome diferente: esperado %q, recebido %q",
			pessoa.Nome, pessoaCriada.Nome)
	}

	if buscarPessoa.EhResponsavel {
		t.Error("esperava pessoa não responsável")
	}

	if buscarPessoa.DataCadastro.IsZero() {
		t.Error("esperava data de cadastro preenchida")
	}

	if buscarPessoa.ID != pessoaCriada.ID {
		t.Errorf(
			"ID diferente: esperado %d, recebido %d",
			pessoaCriada.ID,
			buscarPessoa.ID,
		)
	}

}

func TestPessoaRepository_ListarResponsaveis(t *testing.T) {
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

	repositorio := NovoPessoaRepository(pool)

	pessoaResponsavel1 := domain.Pessoa{
		RegistroPublico: "TEST1-" + rand.Text(),
		Nome:            "Pessoa de teste",
	}
	pessoaResponsavel2 := domain.Pessoa{
		RegistroPublico: "TEST2-" + rand.Text(),
		Nome:            "Pessoa de teste",
	}
	pessoa := domain.Pessoa{
		RegistroPublico: "TEST3-" + rand.Text(),
		Nome:            "Pessoa de teste",
	}

	codigosTeste := []string{
		pessoaResponsavel1.RegistroPublico,
		pessoaResponsavel2.RegistroPublico,
		pessoa.RegistroPublico,
	}

	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM pessoas WHERE registro_publico = ANY($1::text[])",
			codigosTeste,
		)
		if err != nil {
			t.Errorf("não foi possível remover pessoas de teste: %v", err)
		}
	})

	pessoaResponsavel1Criada, err := repositorio.Criar(ctx, pessoaResponsavel1)
	if err != nil {
		t.Fatalf("não foi possível criar pessoa responsavel: %v", err)
	}

	pessoaResponsavel1Criada.EhResponsavel = true
	_, err = repositorio.Atualizar(ctx, pessoaResponsavel1Criada)
	if err != nil {
		t.Fatalf("não foi possível atulizar resposavel: %v", err)
	}

	pessoaResponsavel2Criada, err := repositorio.Criar(ctx, pessoaResponsavel2)
	if err != nil {
		t.Fatalf("não foi possível criar pessoa responsavel: %v", err)
	}

	pessoaResponsavel2Criada.EhResponsavel = true
	_, err = repositorio.Atualizar(ctx, pessoaResponsavel2Criada)
	if err != nil {
		t.Fatalf("não foi possível atulizar resposavel: %v", err)
	}

	pessoaCriada, err := repositorio.Criar(ctx, pessoa)
	if err != nil {
		t.Fatalf("não foi possível criar pessoa: %v", err)
	}

	responsaveis, err := repositorio.ListarResponsaveis(ctx)
	if err != nil {
		t.Fatalf("não foi possível listar responsáveis: %v", err)
	}

	encontrouPrimeiro := false
	encontrouSegundo := false

	for _, responsavel := range responsaveis {
		if !responsavel.EhResponsavel {
			t.Errorf("lista retornou pessoa que não é responsável: id=%d", responsavel.ID)
		}

		if responsavel.ID == pessoaResponsavel1Criada.ID {
			encontrouPrimeiro = true
		}

		if responsavel.ID == pessoaResponsavel2Criada.ID {
			encontrouSegundo = true
		}

		if responsavel.ID == pessoaCriada.ID {
			t.Errorf("lista retornou a pessoa comum: id=%d", pessoaCriada.ID)
		}
	}

	if !encontrouPrimeiro {
		t.Error("primeiro responsável criado não apareceu na lista")
	}

	if !encontrouSegundo {
		t.Error("segundo responsável criado não apareceu na lista")
	}

}

func TestPessoaRepository_Atualizar(t *testing.T) {
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

	repositorio := NovoPessoaRepository(pool)

	pessoa := domain.Pessoa{
		RegistroPublico: "TEST-" + rand.Text(),
		Nome:            "Pessoa de teste",
	}

	t.Cleanup(func() {
		_, err := pool.Exec(
			context.Background(),
			"DELETE FROM pessoas WHERE registro_publico = $1",
			pessoa.RegistroPublico,
		)
		if err != nil {
			t.Errorf("não foi possível remover pessoas de teste: %v", err)
		}
	})

	pessoaCriada, err := repositorio.Criar(ctx, pessoa)
	if err != nil {
		t.Fatalf("não foi possível criar pessoa: %v", err)
	}

	pessoaCriada.Nome = "Nome atualizado"
	pessoaCriada.EhResponsavel = true

	pessoaAtualizada, err := repositorio.Atualizar(ctx, pessoaCriada)
	if err != nil {
		t.Fatalf("não foi possível atualziar pessoa: %v", err)
	}

	if pessoaAtualizada.ID != pessoaCriada.ID {
		t.Errorf("ID diferente: esperado %d, recebido %d",
			pessoaCriada.ID, pessoaAtualizada.ID)
	}

	if pessoaAtualizada.RegistroPublico != pessoaCriada.RegistroPublico {
		t.Errorf(
			"registro público diferente: esperado %q, recebido %q",
			pessoaCriada.RegistroPublico,
			pessoaAtualizada.RegistroPublico,
		)
	}

	if pessoaAtualizada.Nome != "Nome atualizado" {
		t.Errorf(
			"nome diferente: esperado %q, recebido %q",
			"Nome atualizado",
			pessoaAtualizada.Nome,
		)
	}

	if !pessoaAtualizada.EhResponsavel {
		t.Error("esperava pessoa habilitada como responsável")
	}

	if !pessoaAtualizada.DataCadastro.Equal(pessoaCriada.DataCadastro) {
		t.Errorf(
			"data de cadastro foi alterada: esperada %v, recebida %v",
			pessoaCriada.DataCadastro,
			pessoaAtualizada.DataCadastro,
		)
	}

}
