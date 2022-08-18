package links

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LinksRepository struct {
	db     *gorm.DB
	logger *logrus.Entry
	lock   sync.RWMutex
}

func CreateLinksRepository(db *gorm.DB, logger *logrus.Entry) *LinksRepository {
	repo := &LinksRepository{
		db:     db,
		logger: logger,
	}
	return repo.Init()
}

func (r *LinksRepository) Init() *LinksRepository {
	r.db.AutoMigrate(&LinksModel{})
	return r
}

func (r *LinksRepository) Save(link *LinksModel) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.db.Save(link).Error
}

func (r *LinksRepository) GetUnreaden(before time.Duration) ([]*LinksModel, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	var links []*LinksModel
	err := r.db.Order("last_ask asc").Where("last_ask < ?", time.Now().Add(-before)).Find(&links).Error
	return links, err
}