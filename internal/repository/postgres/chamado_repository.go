package postgres

import (
	"context"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChamadoRepository struct {
	pool *pgxpool.Pool
}

func NovoChamadoRepository(pool *pgxpool.Pool) *ChamadoRepository {
	return &ChamadoRepository{
		pool: pool,
	}
}

func (r *ChamadoRepository) Criar(ctx context.Context, chamado domain.Chamado) (domain.Chamado, error) {
	const consulta = `
		INSERT INTO chamados (solicitante_id,responsavel_id,titulo,descricao,prioridade,status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, data_abertura
	`

	err := r.pool.QueryRow(ctx,
		consulta,
		chamado.Solicitante.ID,
		chamado.Responsavel.ID,
		chamado.Titulo,
		chamado.Descricao,
		chamado.Prioridade,
		chamado.Status,
	).Scan(
		&chamado.ID,
		&chamado.DataAbertura,
	)
	if err != nil {
		return domain.Chamado{}, err
	}

	return chamado, nil
}

func (r *ChamadoRepository) BuscarPorID(ctx context.Context, id int64) (domain.Chamado, error) {
	const consulta = `
		SELECT 
			c.id, c.titulo,	c.descricao, 
			c.prioridade, c.status, c.data_abertura,

			s.id, 	s.registro_publico, 	s.nome,
			s.eh_responsavel, 	s.data_cadastro,
			
			r.id, 	r.registro_publico, 	r.nome,
			r.eh_responsavel, 	r.data_cadastro

		FROM chamados AS c
		INNER JOIN pessoas AS s
			ON s.id  = c.solicitante_id
		INNER JOIN  pessoas AS r
			ON r.id = c.responsavel_id
		WHERE c.id = $1
	`
	var chamado domain.Chamado

	err := r.pool.QueryRow(ctx, consulta, id).Scan(
		&chamado.ID,
		&chamado.Titulo,
		&chamado.Descricao,
		&chamado.Prioridade,
		&chamado.Status,
		&chamado.DataAbertura,
		&chamado.Solicitante.ID,
		&chamado.Solicitante.RegistroPublico,
		&chamado.Solicitante.Nome,
		&chamado.Solicitante.EhResponsavel,
		&chamado.Solicitante.DataCadastro,
		&chamado.Responsavel.ID,
		&chamado.Responsavel.RegistroPublico,
		&chamado.Responsavel.Nome,
		&chamado.Responsavel.EhResponsavel,
		&chamado.Responsavel.DataCadastro,
	)

	if err != nil {
		return domain.Chamado{}, err
	}

	return chamado, nil

}

func (r *ChamadoRepository) Listar(ctx context.Context) ([]domain.Chamado, error) {
	const consulta = `
		SELECT 
			id,
			solicitante_id,
			responsavel_id,
			titulo,
			descricao,
			prioridade,
			status,
			data_abertura
		FROM chamados 
		ORDER BY data_abertura DESC, id DESC
	`
	pessoasID := make(map[int64]bool, 0)

	chamados := []domain.Chamado{}

	rows, err := r.pool.Query(ctx, consulta)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var chamado domain.Chamado

		err = rows.Scan(
			&chamado.ID,
			&chamado.Solicitante.ID,
			&chamado.Responsavel.ID,
			&chamado.Titulo,
			&chamado.Descricao,
			&chamado.Prioridade,
			&chamado.Status,
			&chamado.DataAbertura,
		)
		if err != nil {
			return nil, err
		}

		chamados = append(chamados, chamado)
		pessoasID[chamado.Solicitante.ID] = true
		pessoasID[chamado.Responsavel.ID] = true

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()

	ids := make([]int64, 0, len(pessoasID))
	for id := range pessoasID {
		ids = append(ids, id)
	}

	pessoas, err := r.buscarPessoasPorID(ctx, ids)
	if err != nil {
		return []domain.Chamado{}, err
	}

	pessoasPorID := make(map[int64]domain.Pessoa, len(pessoas))
	for _, pessoa := range pessoas {
		pessoasPorID[pessoa.ID] = pessoa
	}

	for i := range chamados {
		chamado := &chamados[i]

		solicitanteID := chamado.Solicitante.ID
		responsavelID := chamado.Responsavel.ID

		chamado.Solicitante = pessoasPorID[solicitanteID]
		chamado.Responsavel = pessoasPorID[responsavelID]
	}

	return chamados, nil

}

func (r *ChamadoRepository) buscarPessoasPorID(ctx context.Context, ids []int64) ([]domain.Pessoa, error) {

	const consulta = `
		SELECT id, registro_publico, nome, eh_responsavel, data_cadastro
		FROM pessoas
		WHERE id = ANY($1::bigint[])
	`

	rows, err := r.pool.Query(ctx, consulta, ids)
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

/*
func(r *ChamadoRepository)Atualizar(ctx context.Context, chamado domain.Chamado) (domain.Chamado, error){

}


for _, pessoa := range pessoas {

			if chamado.Solicitante.ID == pessoa.ID {
				chamado.Solicitante.RegistroPublico = pessoa.RegistroPublico
				chamado.Solicitante.Nome = pessoa.Nome
				chamado.Solicitante.EhResponsavel = pessoa.EhResponsavel
				chamado.Solicitante.DataCadastro = pessoa.DataCadastro
			}

			if chamado.Responsavel.ID == pessoa.ID {
				chamado.Responsavel.RegistroPublico = pessoa.RegistroPublico
				chamado.Responsavel.Nome = pessoa.Nome
				chamado.Responsavel.EhResponsavel = pessoa.EhResponsavel
				chamado.Responsavel.DataCadastro = pessoa.DataCadastro
			}

		}
*/
