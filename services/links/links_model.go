package links

import (
	"time"

	"github.com/guionardo/gs-bot/dal"
)

type (
	LinksModel struct {
		dal.Model
		ChatID         int64  `gorm:"primaryKey;autoIncrement:false"`
		URL            string `gorm:"primaryKey;autoIncrement:false"`
		Title          string `gorm:"type:varchar(200)"`
		Description    string
		Image          string
		SiteName       string
		LastChecked    time.Time
		LastStatusCode int
		LastStatus     string
		LastAsk        time.Time
	}
)
