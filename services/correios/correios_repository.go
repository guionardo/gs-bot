package correios

import (
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CorreiosRepository struct {
	db     *gorm.DB
	logger *logrus.Entry
	lock   sync.RWMutex
}

func CreateCorreiosRepository(db *gorm.DB, logger *logrus.Entry) *CorreiosRepository {
	repo := &CorreiosRepository{
		db:     db,
		logger: logger,
	}
	return repo.Init()
}

func (r *CorreiosRepository) Init() *CorreiosRepository {
	r.db.AutoMigrate(&CorreiosRastreioModel{})
	return r
}

func (r *CorreiosRepository) RastreamentoExiste(codRastreamento string) *CorreiosRastreioModel {
	r.lock.RLock()
	defer r.lock.RUnlock()
	rastreio := &CorreiosRastreioModel{}
	r.db.Where("cod_objeto = ?", codRastreamento).First(rastreio)
	return rastreio
}

func (r *CorreiosRepository) AdicionarRastreamentoAoChat(rastreio *CorreiosRastreioModel, chatID int64) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	return nil
	//TODO: REmover
}

func (r *CorreiosRepository) GetRastreio(codRastreamento string, chatID int64) *CorreiosRastreioModel {
	r.lock.RLock()
	defer r.lock.RUnlock()
	rastreio := &CorreiosRastreioModel{}
	r.db.Where("cod_objeto = ? and chat_id = ?", codRastreamento, chatID).First(rastreio)
	return rastreio
}

func (r *CorreiosRepository) GetRastreios(codRastreamento string) (rastreios []*CorreiosRastreioModel, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	err = r.db.Where("cod_objeto = ?", codRastreamento).Find(&rastreios).Error
	return
}

func (r *CorreiosRepository) GetRastreiosFromChat(chatId int64) (rastreios []*CorreiosRastreioModel, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	err = r.db.Where("chat_id = ?", chatId).Find(&rastreios).Error
	return
}

func (r *CorreiosRepository) RemoveRastreamento(chatID int64, codObjeto string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.db.Where("chat_id = ? and cod_objeto = ?", chatID, codObjeto).Delete(&CorreiosRastreioModel{}).Error
}

func (r *CorreiosRepository) SaveRastreio(rastreio *CorreiosRastreioModel) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.db.Save(rastreio).Error
}

func (r *CorreiosRepository) GetRastreiosPendentes() (rastreios []*CorreiosRastreioModel, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	err = r.db.Where("objeto_entregue = ?", false).Find(&rastreios).Error
	return
}

func (r *CorreiosRepository) GetCodigosRastreioPendentes() (codigos []string, err error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	err = r.db.Table("correios_rastreios").
		Distinct("cod_objeto").
		Select("cod_objeto").
		Where("objeto_entregue = ?", false).
		Pluck("cod_objeto", &codigos).
		Error
	return
}
