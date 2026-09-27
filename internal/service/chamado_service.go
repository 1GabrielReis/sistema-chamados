package service

import (
	"github.com/1GabrielReis/sistema-chamados/internal/repository"
)

type ChamadoService struct {
	repositorioChamado repository.ChamadoRepository
	repositorioPessoa  repository.PessoaRepository
}

func NovoChamadoService(repositorioChamado repository.ChamadoRepository,
	repositorioPessoa repository.PessoaRepository) *ChamadoService {
	return &ChamadoService{
		repositorioChamado: repositorioChamado,
		repositorioPessoa:  repositorioPessoa,
	}
}
