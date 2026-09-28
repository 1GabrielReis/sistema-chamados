package service

import (
	"testing"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
)

func TestValidarDadosChamado(t *testing.T) {
	_, _, err := validarDadosChamado(" ", "uma descrição válida", domain.PrioridadeBaixa)

	if err == nil {
		t.Error("esperava um erro para titulo vazio")
	}

}
