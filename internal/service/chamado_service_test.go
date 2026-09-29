package service

import (
	"testing"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
)

func TestValidarDadosChamado_TituloVazio(t *testing.T) {
	_, _, err := validarDadosChamado(" ", "uma descrição válida", domain.PrioridadeBaixa)

	if err == nil {
		t.Error("esperava um erro para titulo vazio")
	}

}

func TestValidarDadosChamado_DescricaoVazia(t *testing.T) {
	_, _, err := validarDadosChamado("validação", " ", domain.PrioridadeAlta)

	if err == nil {
		t.Error("esperava um erro para descrição vazio")
	}
}

func TestValidarDadosChamado_PrioridadeInvalida(t *testing.T) {
	_, _, err := validarDadosChamado("validação", "uma descrição válida", domain.Prioridade("urgente"))

	if err == nil {
		t.Error("sperava um erro para prioridade inválida")
	}
}

func TestValidarDadosChamado_DadosValidos(t *testing.T) {
	_, _, err := validarDadosChamado("texto válido;", "uma constante válida", domain.PrioridadeAlta)

	if err != nil {
		t.Errorf("não esperava erro para dados válidos: %v", err)
	}
}
