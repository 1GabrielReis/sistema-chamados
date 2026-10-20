package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	var entrada cadastrarPessoaEntrada

	err := json.NewDecoder(r.Body).Decode(&entrada)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	pessoa, err := h.pessoaService.Cadastrar(r.Context(), entrada.Nome)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	saida := pessoaResposta{
		RegistroPublico: pessoa.RegistroPublico,
		Nome:            pessoa.Nome,
		EhResponsavel:   pessoa.EhResponsavel,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(saida); err != nil {
		fmt.Printf("erro ao gerar resposta JSON: %v\n", err)
	}
}
