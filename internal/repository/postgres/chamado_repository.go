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

/*

func (r *ChamadoRepository) buscarPessoaPorID(ctx context.Context, id int64) (domain.Pessoa, error) {
	const cosulta = `
		SELECT id, registro_publico, nome, eh_responsavel, data_cadastro
		FROM pessoas
		WHERE id = $1
	`
	var pessoa domain.Pessoa

	err := r.pool.QueryRow(ctx, cosulta, id).Scan(
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

func(r *ChamadoRepository)Listar(ctx context.Context) ([]domain.Chamado, error){

}
func(r *ChamadoRepository)Atualizar(ctx context.Context, chamado domain.Chamado) (domain.Chamado, error){

}

*/
