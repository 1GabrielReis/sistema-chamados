package handler

type pessoaEntrada struct {
	Nome string `json:"nome"`
}

type pessoaResposta struct {
	RegistroPublico string `json:"registro_publico"`
	Nome            string `json:"nome"`
	EhResponsavel   bool   `json:"eh_responsavel"`
}
