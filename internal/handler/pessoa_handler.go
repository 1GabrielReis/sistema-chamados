package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/1GabrielReis/sistema-chamados/internal/domain"
	"github.com/1GabrielReis/sistema-chamados/internal/service"
)

type PessoaHandler struct {
	pessoaService *service.PessoaService
}

func NovoPessoaHandler(pessoaService *service.PessoaService) *PessoaHandler {
	return &PessoaHandler{
		pessoaService: pessoaService,
	}
}

func (h *PessoaHandler) Cadastrar(w http.ResponseWriter, r *http.Request) {
	var entrada pessoaEntrada

	err := json.NewDecoder(r.Body).Decode(&entrada)
	if err != nil {
		fmt.Printf("erro ao interpretar JSON do cadastro de pessoa: %v\n", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	pessoa, err := h.pessoaService.Cadastrar(r.Context(), entrada.Nome)
	if err != nil {
		fmt.Printf("erro ao cadastrar pessoa: %v\n", err)
		http.Error(w, "não foi possível cadastrar pessoa", http.StatusBadRequest)
		return
	}

	saida := novaPessoaResposta(pessoa)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON de pessoa: %v\n", err)
	}
}

func (h *PessoaHandler) BuscarPorCodigoPublico(w http.ResponseWriter, r *http.Request) {

	codigo := r.PathValue("codigo")

	pessoa, err := h.pessoaService.BuscarPorCodigoPublico(r.Context(), codigo)
	if err != nil {
		fmt.Printf("erro ao buscar pessoa pelo código %q: %v\n", codigo, err)
		http.Error(w, "pessoa não encontrada", http.StatusNotFound)
		return
	}

	saida := novaPessoaResposta(pessoa)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON de pessoa: %v\n", err)
	}

}

func (h *PessoaHandler) ListarResponsaveis(w http.ResponseWriter, r *http.Request) {

	pessoas, err := h.pessoaService.ListarResponsaveis(r.Context())
	if err != nil {
		fmt.Printf("erro ao listar responsáveis: %v\n", err)
		http.Error(w, "erro interno ao listar responsáveis", http.StatusInternalServerError)
		return
	}

	saida := make([]pessoaResposta, 0, len(pessoas))

	for _, pessoa := range pessoas {
		saida = append(saida, novaPessoaResposta(pessoa))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON de pessoa: %v\n", err)
	}
}

func (h *PessoaHandler) ListarSolicitantes(w http.ResponseWriter, r *http.Request) {
	pessoas, err := h.pessoaService.ListarSolicitantes(r.Context())
	if err != nil {
		fmt.Printf("erro ao listar solicitantes: %v\n", err)
		http.Error(w, "erro interno ao listar solicitantes", http.StatusInternalServerError)
		return
	}

	saida := make([]pessoaResposta, 0, len(pessoas))

	for _, pessoa := range pessoas {
		saida = append(saida, novaPessoaResposta(pessoa))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON de pessoa: %v\n", err)
	}

}

func (h *PessoaHandler) Atualizar(w http.ResponseWriter, r *http.Request) {

	codigo := r.PathValue("codigo")
	var entrada pessoaEntrada

	err := json.NewDecoder(r.Body).Decode(&entrada)
	if err != nil {
		fmt.Printf("erro ao interpretar JSON de atualização de pessoa: %v\n", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	pessoa, err := h.pessoaService.Atualizar(r.Context(), codigo, entrada.Nome)
	if err != nil {
		fmt.Printf("erro ao atualizar pessoa com código %q: %v\n", codigo, err)
		http.Error(w, "não foi possível atualizar pessoa", http.StatusBadRequest)
		return
	}

	saida := novaPessoaResposta(pessoa)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON de pessoa: %v\n", err)
	}
}

func (h *PessoaHandler) HabilitarComoResponsavel(w http.ResponseWriter, r *http.Request) {
	codigo := r.PathValue("codigo")
	var entrada habilitarResponsavelEntrada

	err := json.NewDecoder(r.Body).Decode(&entrada)
	if err != nil {
		fmt.Printf("erro ao interpretar JSON para habilitar pessoa como responsável: %v\n", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	pessoa, err := h.pessoaService.HabilitarComoResponsavel(r.Context(), entrada.CodigoAutorizador, codigo)
	if err != nil {
		fmt.Printf("erro ao habilitar pessoa como responsável: %v\n", err)
		http.Error(w, "não foi possível habilitar pessoa como responsável", http.StatusBadRequest)
		return
	}

	saida := novaPessoaResposta(pessoa)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON de pessoa: %v\n", err)
	}

}

func novaPessoaResposta(pessoa domain.Pessoa) pessoaResposta {
	return pessoaResposta{
		RegistroPublico: pessoa.RegistroPublico,
		Nome:            pessoa.Nome,
		EhResponsavel:   pessoa.EhResponsavel,
	}
}
