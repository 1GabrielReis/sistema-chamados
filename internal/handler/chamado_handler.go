package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
		fmt.Printf("JSON inválido: %v\n", err)
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
		fmt.Printf("erro ao gerar resposta JSON do chamado: %v\n", err)
	}
}

func (h *ChamadoHandler) BuscarPorID(w http.ResponseWriter, r *http.Request) {

	idTexto := r.PathValue("id")

	id, err := strconv.ParseInt(idTexto, 10, 64)
	if err != nil {
		fmt.Printf("ID malformado ou não positivo: %v\n", err)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	if id <= 0 {
		fmt.Printf("ID malformado ou não positivo: %d\n", id)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	chamado, err := h.chamadoService.BuscarPorID(r.Context(), id)
	if err != nil {
		fmt.Printf("erro ao buscar chamado pelo id %d: %v\n", id, err)
		http.Error(w, "chamado não encontrado", http.StatusNotFound)
		return
	}

	saida := novaChamadoResposta(chamado)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON do chamado: %v\n", err)
	}

}

func (h *ChamadoHandler) Listar(w http.ResponseWriter, r *http.Request) {

	chamados, err := h.chamadoService.Listar(r.Context())
	if err != nil {
		fmt.Printf("erro ao listar chamados: %v\n", err)
		http.Error(w, "erro interno ao listar chamados", http.StatusInternalServerError)
		return
	}

	saida := make([]chamadoResposta, 0, len(chamados))

	for _, chamado := range chamados {
		saida = append(saida, novaChamadoResposta(chamado))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON do chamado: %v\n", err)
	}
}

func (h *ChamadoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		fmt.Printf("ID malformado ou não positivoo: %v\n", err)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	if id <= 0 {
		fmt.Printf("ID malformado ou não positivo: %d\n", id)
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var entrada chamadoAtualizacaoEntrada

	err = json.NewDecoder(r.Body).Decode(&entrada)
	if err != nil {
		fmt.Printf("erro ao interpretar JSON de atualização do chamado: %v\n", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	chamado, err := h.chamadoService.Atualizar(r.Context(), novaChamadaAtualizarEntrada(id, entrada))
	if err != nil {
		fmt.Printf("erro ao atualizar chamado de ID %d: %v\n", id, err)
		http.Error(w, "não foi possível atualizar chamado", http.StatusBadRequest)
		return
	}

	saida := novaChamadoResposta(chamado)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON do chamado: %v\n", err)
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

func novaChamadaAtualizarEntrada(id int64, chamadoEntrada chamadoAtualizacaoEntrada) service.AtualizarChamadoEntrada {
	return service.AtualizarChamadoEntrada{
		ID:                     id,
		CodigoResponsavel:      chamadoEntrada.CodigoResponsavel,
		DistribuicaoAutomatica: chamadoEntrada.DistribuicaoAutomatica,
		Titulo:                 chamadoEntrada.Titulo,
		Descricao:              chamadoEntrada.Descricao,
		Prioridade:             domain.Prioridade(chamadoEntrada.Prioridade),
		Status:                 domain.Status(chamadoEntrada.Status),
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
