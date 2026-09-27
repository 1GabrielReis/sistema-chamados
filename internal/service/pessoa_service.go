package service

import (
	"context"
	"crypto/rand"
	"errors"
	"strings"
	"time"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
	"github.com/1GabrielReis/sistema-chamados/internal/repository"
)

type PessoaService struct {
	repositorio repository.PessoaRepository
}

func NovoPessoaService(repositorio repository.PessoaRepository) *PessoaService {
	return &PessoaService{
		repositorio: repositorio,
	}
}

func (s *PessoaService) Cadastrar(ctx context.Context,
	nome string) (domain.Pessoa, error) {
	nome = strings.TrimSpace(nome)
	if nome == "" {
		return domain.Pessoa{}, errors.New("nome é obrigatorio")
	}
	pessoa := domain.Pessoa{
		RegistroPublico: rand.Text(),
		Nome:            nome,
		EhResponsavel:   false,
		DataCadastro:    time.Now(),
	}
	return s.repositorio.Criar(ctx, pessoa)
}

func (s *PessoaService) BuscarPorCodigoPublico(ctx context.Context,
	registroPublico string) (domain.Pessoa, error) {
	registroPublico = strings.TrimSpace(registroPublico)
	if registroPublico == "" {
		return domain.Pessoa{}, errors.New("registro público é obrigatório")
	}
	return s.repositorio.BuscarPorCodigoPublico(ctx, registroPublico)
}

/*
func (s *PessoaService) ListarResponsaveis() {
	s.repositorio.ListarResponsaveis()
}

func (s *PessoaService) Atualizar() {
	s.repositorio.Atualizar()
}
*/
