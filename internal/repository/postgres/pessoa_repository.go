package postgres

import (
	"context"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
	"github.com/1GabrielReis/sistema-chamados/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PessoaRepository struct {
	pool *pgxpool.Pool
}

var _ repository.PessoaRepository = (*PessoaRepository)(nil)

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

func (r *PessoaRepository) Atualizar(ctx context.Context, pessoa domain.Pessoa) (domain.Pessoa, error) {
	const consulta = `
		UPDATE pessoas
		SET  nome = $1, eh_responsavel = $2
		WHERE id = $3 
		RETURNING  RETURNING id, registro_publico, nome, eh_responsavel, data_cadastro
	`
	err := r.pool.QueryRow(ctx,
		consulta, pessoa.Nome,
		pessoa.EhResponsavel,
		pessoa.ID,
	).Scan(
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

func (r *PessoaRepository) BuscarResponsavelComMenosChamadosAbertos(ctx context.Context) (domain.Pessoa, error) {
	const consulta = `
		SELECT
			p.id,
			p.registro_publico,
			p.nome,
			p.eh_responsavel,
			p.data_cadastro
		FROM pessoas AS p
		LEFT JOIN chamados AS c 
			ON c.responsavel_id = p.id
			AND c.status IN ('aberto', 'em_andamento')
		WHERE p.eh_responsavel = TRUE
		GROUP BY
			p.id,
			p.registro_publico,
			p.nome,
			p.eh_responsavel,
			p.data_cadastro
		ORDER BY COUNT(c.id) ASC, p.id ASC
		LIMIT 1
	`
	var pessoa domain.Pessoa

	err := r.pool.QueryRow(ctx, consulta).Scan(
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
