package handler

import "time"

type chamadoCadastroEntrada struct {
	CodigoSolicitante      string `json:"codigo_solicitante"`
	CodigoResponsavel      string `json:"codigo_responsavel"`
	DistribuicaoAutomatica bool   `json:"distribuicao_automatica"`
	Titulo                 string `json:"titulo"`
	Descricao              string `json:"descricao"`
	Prioridade             string `json:"prioridade"`
}

type chamadoResposta struct {
	ID           int64          `json:"id"`
	Solicitante  pessoaResposta `json:"solicitante"`
	Responsavel  pessoaResposta `json:"responsavel"`
	Titulo       string         `json:"titulo"`
	Descricao    string         `json:"descricao"`
	Prioridade   string         `json:"prioridade"`
	Status       string         `json:"status"`
	DataAbertura time.Time      `json:"data_abertura"`
}
