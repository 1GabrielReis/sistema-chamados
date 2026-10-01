package postgres

import (
	"context"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PessoaRepository struct {
	pool *pgxpool.Pool
}

func NovoPessoaRepository(pool *pgxpool.Pool) *PessoaRepository {
	return &PessoaRepository{
		pool: pool,
	}
}

func (r *PessoaRepository) Criar(ctx context.Context, pessoa domain.Pessoa) (domain.Pessoa, error) {
	const consulta = `
					INSERT INTO pessoas (registro_publico,nome) 
					VALUES ($1,$2) 
					RETURNING id, eh_responsavel, data_cadastro
				`

	err := r.pool.QueryRow(ctx,
		consulta,
		pessoa.RegistroPublico,
		pessoa.Nome,
	).Scan(
		&pessoa.ID,
		&pessoa.EhResponsavel,
		&pessoa.DataCadastro,
	)

	if err != nil {
		return domain.Pessoa{}, err
	}

	return pessoa, err

}

func (r *PessoaRepository) BuscarPorCodigoPublico(ctx context.Context, registroPublico string) (domain.Pessoa, error) {
	const consulta = `
		SELECT id, registro_publico, nome, eh_responsavel, data_cadastro
		FROM pessoas 
		WHERE registro_publico = $1
	`
	var pessoa domain.Pessoa

	err := r.pool.QueryRow(ctx, consulta, registroPublico).Scan(
		&pessoa.ID,
		&pessoa.RegistroPublico,
		&pessoa.Nome,
		&pessoa.EhResponsavel,
		&pessoa.DataCadastro,
	)

	if err != nil {
		return domain.Pessoa{}, err
	}

	return pessoa, nil
}
