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
