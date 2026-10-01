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

func (r *PessoaRepository) ListarResponsaveis(ctx context.Context) ([]domain.Pessoa, error) {
	const consulta = `
		SELECT id, registro_publico, nome, eh_responsavel, data_cadastro
		FROM pessoas 
		WHERE eh_responsavel = TRUE
		ORDER BY nome
	`
	rows, err := r.pool.Query(ctx, consulta)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pessoas := make([]domain.Pessoa, 0)

	for rows.Next() {
		var pessoa domain.Pessoa

		err = rows.Scan(
			&pessoa.ID,
			&pessoa.RegistroPublico,
			&pessoa.Nome,
			&pessoa.EhResponsavel,
			&pessoa.DataCadastro,
		)
		if err != nil {
			return nil, err
		}

		pessoas = append(pessoas, pessoa)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pessoas, nil

}
