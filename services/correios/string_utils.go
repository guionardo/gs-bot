package correios

import (
	correio "github.com/guionardo/go-gstools/correios"
	"github.com/guionardo/gs-bot/internal"
)
func unidadeString(um correio.Unidade) string {
	return internal.JoinNotEmpty(um.Tipo, um.Nome, um.CodSro, um.Endereco.String())
}