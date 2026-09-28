package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
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

type CadastrarChamadoEntrada struct {
	CodigoSolicitante      string
	CodigoResponsavel      string
	DistribuicaoAutomatica bool
	Titulo                 string
	Descricao              string
	Prioridade             domain.Prioridade
}

type AtualizarChamadoEntrada struct {
	ID                int64
	CodigoResponsavel string
	Titulo            string
	Descricao         string
	Prioridade        domain.Prioridade
	Status            domain.Status
}

func (s *ChamadoService) Cadastrar(ctx context.Context,
	entrada CadastrarChamadoEntrada) (domain.Chamado, error) {
	titulo, descricao, err := validarDadosChamado(
		entrada.Titulo,
		entrada.Descricao,
		entrada.Prioridade,
	)
	if err != nil {
		return domain.Chamado{}, err
	}

	entrada.CodigoSolicitante = strings.TrimSpace(entrada.CodigoSolicitante)
	if entrada.CodigoSolicitante == "" {
		return domain.Chamado{}, errors.New("código do responsável é obrigatório")
	}

	solicitante, err := s.repositorioPessoa.BuscarPorCodigoPublico(ctx, entrada.CodigoSolicitante)
	if err != nil {
		return domain.Chamado{}, err
	}

	var responsavel domain.Pessoa

	if !entrada.DistribuicaoAutomatica {
		entrada.CodigoResponsavel = strings.TrimSpace(entrada.CodigoResponsavel)
		if entrada.CodigoResponsavel == "" {
			return domain.Chamado{}, errors.New("código do solicitante é obrigatório")
		}
		responsavel, err = s.repositorioPessoa.BuscarPorCodigoPublico(ctx, entrada.CodigoResponsavel)
		if err != nil {
			return domain.Chamado{}, err
		}
	} else {
		responsavel, err = s.repositorioPessoa.BuscarResponsavelComMenosChamadosAbertos(ctx)
		if err != nil {
			return domain.Chamado{}, err
		}
	}

	if !responsavel.EhResponsavel {
		return domain.Chamado{}, errors.New("pessoa informada não é responsável")
	}

	chamado := domain.Chamado{
		Solicitante:  solicitante,
		Responsavel:  responsavel,
		Titulo:       titulo,
		Descricao:    descricao,
		Prioridade:   entrada.Prioridade,
		Status:       domain.StatusAberto,
		DataAbertura: time.Now(),
	}
	return s.repositorioChamado.Criar(ctx, chamado)
}

func (s *ChamadoService) BuscarPorID(ctx context.Context,
	id int64) (domain.Chamado, error) {
	if id <= 0 {
		return domain.Chamado{}, errors.New("id do chamado inválido")
	}
	return s.repositorioChamado.BuscarPorID(ctx, id)
}

func (s *ChamadoService) Listar(ctx context.Context) ([]domain.Chamado, error) {
	return s.repositorioChamado.Listar(ctx)
}

func validarDadosChamado(titulo string, descricao string,
	prioridade domain.Prioridade) (string, string, error) {
	titulo = strings.TrimSpace(titulo)
	if titulo == "" {
		return "", "", errors.New("título é obrigatório")
	}

	descricao = strings.TrimSpace(descricao)
	if descricao == "" {
		return "", "", errors.New("descrição é obrigatória")
	}

	switch prioridade {
	case domain.PrioridadeBaixa,
		domain.PrioridadeMedia,
		domain.PrioridadeAlta:
		// Prioridade válida; continua o método.
	default:
		return "", "", errors.New("prioridade inválida")
	}
	return titulo, descricao, nil
}
