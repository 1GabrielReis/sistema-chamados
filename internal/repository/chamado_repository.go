package repository

import (
	"context"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
)

type ChamadoRepository interface {
	Criar(ctx context.Context, chamado domain.Chamado) (domain.Chamado, error)
	BuscarPorID(ctx context.Context, id int64) (domain.Chamado, error)
	Listar(ctx context.Context) ([]domain.Chamado, error)
	Atualizar(ctx context.Context, chamado domain.Chamado) (domain.Chamado, error)
}
