package repository

import (
	"context"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
)

type PessoaRepository interface {
	Criar(ctx context.Context, pessoa domain.Pessoa) (domain.Pessoa, error)
	BuscarPorCodigoPublico(ctx context.Context, RegistroPublico string) (domain.Pessoa, error)
	ListarResponsaveis(ctx context.Context) ([]domain.Pessoa, error)
	Atualizar(ctx context.Context, pessoa domain.Pessoa) (domain.Pessoa, error)
	BuscarResponsavelComMenosChamadosAbertos(ctx context.Context) (domain.Pessoa, error)
}
