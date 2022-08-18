package correios

import (
	"encoding/json"

	"github.com/guionardo/gs-bot/internal"
)

type CorreiosResponseModel struct {
	Objetos    []CorreiosObjetoModel `json:"objetos"`
	Quantidade int                   `json:"quantidade"`
	Resultado  string                `json:"resultado"`
	Versao     string                `json:"versao"`
}

func (crm *CorreiosResponseModel) Valido() bool {
	return crm.Resultado == "OK" &&
		crm.Quantidade > 0 &&
		crm.Objetos != nil &&
		len(crm.Objetos[0].Eventos) > 0
}

type CorreiosObjetoModel struct {
	CodObjeto  string                  `json:"codObjeto"`
	Eventos    []CorreiosEventoModel   `json:"eventos"`
	Modalidade string                  `json:"modalidade"`
	TipoPostal CorreiosTipoPostalModel `json:"tipoPostal"`
}

type CorreiosEventoModel struct {
	Codigo         string                     `json:"codigo"`
	Descricao      string                     `json:"descricao"`
	DtHrCriado     string                     `json:"dtHrCriado"`
	Tipo           string                     `json:"tipo"`
	Unidade        CorreiosEventoUnidadeModel `json:"unidade"`
	UnidadeDestino CorreiosEventoUnidadeModel `json:"unidadeDestino"`
}

type CorreiosEventoUnidadeModel struct {
	CodSro   string                             `json:"codSro"`
	Endereco CorreiosEventoUnidadeEnderecoModel `json:"endereco"`
	Nome     string                             `json:"nome"`
	Tipo     string                             `json:"tipo"`
}

func (um *CorreiosEventoUnidadeModel) String() string {
	return internal.JoinNotEmpty(um.Tipo, um.Nome, um.CodSro, um.Endereco.String())
}

type CorreiosEventoUnidadeEnderecoModel struct {
	Bairro     string `json:"bairro"`
	CEP        string `json:"cep"`
	Cidade     string `json:"cidade"`
	Logradouro string `json:"logradouro"`
	Numero     string `json:"numero"`
	UF         string `json:"uf"`
}

func (em *CorreiosEventoUnidadeEnderecoModel) String() string {
	return internal.JoinNotEmpty(em.Logradouro, em.Numero, em.Bairro, em.Cidade, em.UF)
}

type CorreiosTipoPostalModel struct {
	Categoria string `json:"categoria"`
	Descricao string `json:"descricao"`
	Sigla     string `json:"sigla"`
}

func ParseResponseModel(body []byte) (crm *CorreiosResponseModel, err error) {
	crm = &CorreiosResponseModel{}
	err = json.Unmarshal(body, crm)
	return
}