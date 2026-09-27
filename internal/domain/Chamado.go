package domain

import "time"

type Prioridade string
type Status string

type Chamado struct {
	ID           int64
	Solicitante  Pessoa
	Responsavel  Pessoa
	Titulo       string
	Descricao    string
	Prioridade   Prioridade
	Status       Status
	DataAbertura time.Time
}

const (
	PrioridadeBaixa Prioridade = "baixa"
	PrioridadeMedia Prioridade = "media"
	PrioridadeAlta  Prioridade = "alta"
)

const (
	StatusAberto      Status = "aberto"
	StatusEmAndamento Status = "em_andamento"
	StatusResolvido   Status = "resolvido"
	StatusFechado     Status = "fechado"
)
