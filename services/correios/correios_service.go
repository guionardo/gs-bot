package correios

import (
	"fmt"

	correio "github.com/guionardo/go-gstools/correios"
	"github.com/sirupsen/logrus"
)

type CorreiosService struct {
	repository *CorreiosRepository
	logger     *logrus.Entry
}

func (svc *CorreiosService) AdicionarRastreamento(codRastreamento string, chatID int64) (err error) {
	rastreio := svc.repository.RastreamentoExiste(codRastreamento)
	if rastreio != nil {
		err = svc.repository.AdicionarRastreamentoAoChat(rastreio, chatID)
		return fmt.Errorf("Rastreamento já existe: %s", rastreio.String())
	}

	rastreamento, err := correio.GetRastreio(codRastreamento)
	if err == nil && !rastreamento.Valido() {
		err = fmt.Errorf("Rastreamento não encontrado: %s", codRastreamento)
	}

	if err != nil {
		return err
	}
	rastreio = GetRastreioModel(rastreamento, chatID)
	err = svc.repository.SaveRastreio(rastreio)

	if err != nil {
		return err
	}
	return nil
}

func (svc *CorreiosService) VerificarRastreamento(codigoRastreamento string, chatID int64) (atualizado bool, ultimoEstado *CorreiosRastreioModel, err error) {
	estadoAnterior := svc.repository.GetRastreio(codigoRastreamento, chatID)
	rastreamento, err := correio.GetRastreio(codigoRastreamento)
	if err != nil {
		return false, nil, err
	}
	rastreio := GetRastreioModel(rastreamento, chatID)
	if rastreio == nil {
		return false, nil, fmt.Errorf("Rastreamento não encontrado: %s", codigoRastreamento)
	}
	if estadoAnterior != nil && rastreio.DataHoraEvento.Equal(estadoAnterior.DataHoraEvento) {
		return false, rastreio, nil
	}
	err = svc.repository.SaveRastreio(rastreio)
	if err != nil {
		return false, rastreio, err
	}
	return true, rastreio, nil
}
