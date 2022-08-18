package correios

import (
	"reflect"
	"testing"

	"github.com/guionardo/go-tgbot/tgbot/config"
	"github.com/guionardo/go-tgbot/tgbot/infra"
	"github.com/guionardo/gs-bot/configuration"
	"github.com/guionardo/gs-bot/dal"
	"github.com/guionardo/gs-bot/internal"
	"github.com/sirupsen/logrus"
)

func TestCorreiosService_Rastreio(t *testing.T) {
	logCfg :=&config.LoggerConfiguration{
		Level: "debug",
	}
	infra.CreateLoggerFactory(logCfg)
		
	cfg := configuration.RepositoryConfiguration{
		ConnectionString: "test.db",
	}
	db, err := dal.GetDatabase(cfg)
	if err != nil {
		t.Error(err)
	}
	repository := CreateCorreiosRepository(db, infra.GetLogger("correios_repository")).Init()

	type fields struct {
		repository *CorreiosRepository
		logger     *logrus.Entry
	}
	type args struct {
		codigoRastreamento string
	}

	tests := []struct {
		name         string
		fields       fields
		args         args
		wantRastreio *CorreiosResponseModel
		wantErr      bool
	}{
		{
			name: "Rastreio de encomenda correios",
			fields: fields{
				repository: repository,
				logger:     internal.GetLogger().WithField("service", "correios"),
			},
			args: args{
				codigoRastreamento: "LB330827204HK",
			},
			wantRastreio: &CorreiosResponseModel{},
			wantErr:      false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &CorreiosService{
				repository: tt.fields.repository,
				logger:     tt.fields.logger,
			}
			gotRastreio, _, err := svc.VerificarRastreamento(tt.args.codigoRastreamento, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("CorreiosService.Rastreio() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRastreio, tt.wantRastreio) {
				t.Errorf("CorreiosService.Rastreio() = %v, want %v", gotRastreio, tt.wantRastreio)
			}
		})
	}
}
