package correios

import (
	"fmt"
	"time"

	correio "github.com/guionardo/go-gstools/correios"
	"github.com/guionardo/gs-bot/dal"
)

type (
	CorreiosRastreioModel struct {
		dal.Model
		ChatID         int64  `gorm:"primaryKey;autoIncrement:false"`
		CodObjeto      string `gorm:"primaryKey;autoIncrement:false"`
		DataHoraEvento time.Time
		DescEvento     string `gorm:"type:varchar(100)"`
		Unidade        string `gorm:"type:varchar(200)"`
		UnidadeDestino string `gorm:"type:varchar(200)"`
		ObjetoEntregue bool   `gorm:"type:boolean"`
		Descricao      string `gorm:"type:varchar(200)"`
		IconStatus     string `gorm:"type:varchar(1)"`
	}
)

func (r *CorreiosRastreioModel) String() string {
	descricao := ""
	if len(r.Descricao) > 0 {
		descricao = fmt.Sprintf("[%s] ", r.Descricao)
	}
	return fmt.Sprintf("%s%s%s - %v - %s - %s", r.IconStatus, descricao, r.CodObjeto, r.DataHoraEvento, r.DescEvento, r.Unidade)
}

func GetRastreioModel(rastreamento *correio.Rastreio, chatID int64) *CorreiosRastreioModel {
	if !(rastreamento != nil && len(rastreamento.Objetos) > 0) {
		return nil
	}
	objeto := rastreamento.Objetos[0]
	if objeto.Eventos == nil || len(objeto.Eventos) == 0 {
		return nil
	}
	evento := objeto.Eventos[0]
	dataHoraEvento := time.Time(evento.DataCriado)

	status, ok := correio.Statuses[evento.Codigo]
	if !ok {
		status = correio.Status{			
			Icon: "",
		}
	}
	rastreio := &CorreiosRastreioModel{
		ChatID:         chatID,
		CodObjeto:      objeto.CodObjeto,
		DataHoraEvento: dataHoraEvento,
		DescEvento:     evento.Descricao,
		Unidade:        evento.Unidade.String(),
		UnidadeDestino: evento.UnidadeDestino.String(),
		ObjetoEntregue: evento.Codigo == "BDE",
		IconStatus:     status.Icon,
	}
	return rastreio
}
