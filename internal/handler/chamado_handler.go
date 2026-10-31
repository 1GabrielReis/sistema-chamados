package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
	"github.com/1GabrielReis/sistema-chamados/internal/service"
)

type ChamadoHandler struct {
	chamadoService *service.ChamadoService
}

func NovoChamadoHandler(chamadoService *service.ChamadoService) *ChamadoHandler {
	return &ChamadoHandler{
		chamadoService: chamadoService,
	}
}

func (h *ChamadoHandler) Cadastrar(w http.ResponseWriter, r *http.Request) {
	var entrada chamadoCadastroEntrada

	err := json.NewDecoder(r.Body).Decode(&entrada)
	if err != nil {
		fmt.Printf("erro ao interpretar JSON do cadastro de chamado: %v\n", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	chamado, err := h.chamadoService.Cadastrar(r.Context(), novaChamadaCadastroEntrada(entrada))
	if err != nil {
		fmt.Printf("erro ao cadastrar chamado: %v\n", err)
		http.Error(w, "não foi possível cadastrar chamado", http.StatusBadRequest)
		return
	}

	saida := novaChamadoResposta(chamado)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON: %v\n", err)
	}
}

func novaChamadaCadastroEntrada(chamadoEntrada chamadoCadastroEntrada) service.CadastrarChamadoEntrada {
	return service.CadastrarChamadoEntrada{
		CodigoSolicitante:      chamadoEntrada.CodigoSolicitante,
		CodigoResponsavel:      chamadoEntrada.CodigoResponsavel,
		DistribuicaoAutomatica: chamadoEntrada.DistribuicaoAutomatica,
		Titulo:                 chamadoEntrada.Titulo,
		Descricao:              chamadoEntrada.Descricao,
		Prioridade:             domain.Prioridade(chamadoEntrada.Prioridade),
	}
}

func novaChamadoResposta(chamado domain.Chamado) chamadoResposta {
	return chamadoResposta{
		ID:           chamado.ID,
		Solicitante:  novaPessoaResposta(chamado.Solicitante),
		Responsavel:  novaPessoaResposta(chamado.Responsavel),
		Descricao:    chamado.Descricao,
		Titulo:       chamado.Titulo,
		Prioridade:   string(chamado.Prioridade),
		Status:       string(chamado.Status),
		DataAbertura: chamado.DataAbertura,
	}
}
