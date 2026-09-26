package domain

import "time"

type Pessoa struct {
	ID              int64
	RegistroPublico string
	Nome            string
	EhResponsavel   bool
	DataCadastro    time.Time
}
