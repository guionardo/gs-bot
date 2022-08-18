package correios

import (
	correio "github.com/guionardo/go-gstools/correios"
	"github.com/guionardo/go-tgbot/tgbot/infra"
)

type RastreioEntidade struct {
	Model    *CorreiosRastreioModel
	Response *correio.Rastreio
}

func GetRastreiosPendentes() (rastreios []*RastreioEntidade, err error) {
	rastreiosDB, err := correiosService.repository.GetRastreiosPendentes()
	if err != nil {
		return nil, err
	}
	rastreios = make([]*RastreioEntidade, len(rastreiosDB))
	for index, rastreioDB := range rastreiosDB {
		rastreioEntidade := &RastreioEntidade{
			Model: rastreioDB,
		}
		rastreios[index] = rastreioEntidade
	}
	return
}

func GetRastreioEntidades(codigoRastreio string) (rastreios []RastreioEntidade, err error) {
	rastreiosDB, err := correiosService.repository.GetRastreios(codigoRastreio)
	logger := infra.GetLogger("correios")
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	rastreios = make([]RastreioEntidade, len(rastreiosDB))
	rastreio, err := correio.GetRastreio(codigoRastreio)
	if err != nil {
		logger.Warningf("Erro ao buscar rastreio: %s", err)
	}
	for index, rastreioDB := range rastreiosDB {
		rastreioEntidade := RastreioEntidade{
			Model:    rastreioDB,
			Response: rastreio,
		}
		rastreios[index] = rastreioEntidade
	}
	return
}

// Atualiza entidade no banco de dados a partir da resposta da API dos correios
func (re *RastreioEntidade) Update() (atualizado bool) {
	atualizado = false
	if re.Response == nil ||
		re.Response.Objetos == nil ||
		len(re.Response.Objetos) == 0 ||
		re.Response.Objetos[0].Eventos == nil ||
		len(re.Response.Objetos[0].Eventos) == 0 {
		// Não há rastreamento nos correios
		return
	}
	chatID := int64(0)
	if re.Model == nil {
		re.Model = &CorreiosRastreioModel{
			ChatID:    chatID,
			CodObjeto: re.Response.Objetos[0].CodObjeto,
		}
	} else {
		chatID = re.Model.ChatID
	}

	modelEvento := GetRastreioModel(re.Response, chatID)

	if !re.Model.DataHoraEvento.Equal(modelEvento.DataHoraEvento) ||
		re.Model.ObjetoEntregue != modelEvento.ObjetoEntregue {
		re.Model.DataHoraEvento = modelEvento.DataHoraEvento
		re.Model.DescEvento = modelEvento.DescEvento
		re.Model.Unidade = modelEvento.Unidade
		re.Model.UnidadeDestino = modelEvento.UnidadeDestino
		re.Model.ObjetoEntregue = modelEvento.ObjetoEntregue
		logger := infra.GetLogger("correios")
		logger.Infof("Rastreio atualizado: %s", re.Model.String())
		atualizado = true
	}
	return
}
