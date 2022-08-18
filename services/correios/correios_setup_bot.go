/*
Comandos

/correio <codigo>	-> Adicionar rastreamento para código
/correio
*/
package correios

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/guionardo/go-tgbot/pkg/schedules"
	"github.com/guionardo/go-tgbot/tgbot"
	"github.com/guionardo/go-tgbot/tgbot/helpers"
	"github.com/guionardo/go-tgbot/tgbot/infra"
	"github.com/guionardo/gs-bot/configuration"
	"github.com/guionardo/gs-bot/dal"
)

var (
	correiosService    *CorreiosService
	correiosRepository *CorreiosRepository
)

func SetupBot(svc *tgbot.GoTGBotService) {
	cfg, err := configuration.GetConfiguration()
	if err != nil {
		panic(err)
	}
	db, err := dal.GetDatabase(cfg.Repository)
	if err != nil {
		panic(err)
	}
	correiosRepository = CreateCorreiosRepository(db, infra.GetLogger("correios_repo")).Init()
	correiosService = &CorreiosService{
		repository: correiosRepository,
		logger:     infra.GetLogger("correios_service"),
	}
	svc.AddCommandHandlers(
		&tgbot.ListenerCommandHandler{
			Command: "correio",
			Title:   "Rastrear encomenda do Correio",
			Func:    correiosRastrearEncomenda,
		},
	)
	svc.AddCallbackHandlers(tgbot.CreateListenerCallbackHandler("Correios", "correio", correiosCallbackHandler))

	// svc.AddHandlers(&tgbot.ListenerFilteredHandler{
	// 	Title:  "all",
	// 	Filter: func(update tgbotapi.Update) bool { return true },
	// 	Func: func(ctx context.Context, update tgbotapi.Update) error {
	// 		svc := tgbot.GetBotService(ctx)
	// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, "+update.Message.From.UserName+"!")
	// 		msg.ReplyToMessageID = update.Message.MessageID
	// 		svc.Publisher().Publish(msg)
	// 		return nil
	// 	}})

	svc.SetupBackgroundSchedules(schedules.CreateSchedule("Verificar rastreamentos", time.Hour, correiosVerificarRastreios))
}

func correiosCallbackHandler(ctx context.Context, update tgbotapi.Update) error {
	svc := tgbot.GetBotService(ctx)
	chatID := update.CallbackQuery.Message.Chat.ID
	switch update.CallbackQuery.Data {
	case "correio|listar":
		lista, err := correiosService.repository.GetRastreiosFromChat(chatID)
		if err != nil {
			svc.Publisher().SendHTMLMessage(chatID, fmt.Sprintf("%v", err))
			return err
		}
		if len(lista) == 0 {
			svc.Publisher().SendTextMessage(chatID, "Nenhum rastreamento adicionado")
			return nil
		}

		//svc.Publisher().ReplyToMessage(update, "Rastreamentos:")
		content := ""
		for _, rastreio := range lista {
			content += fmt.Sprintf("%s\n", rastreio.String())
		}
		svc.Publisher().SendHTMLMessage(chatID, content)
		return nil
	case "correio|remover":
		lista, err := correiosService.repository.GetRastreiosFromChat(chatID)
		if err != nil {
			svc.Publisher().SendHTMLMessage(chatID, fmt.Sprintf("%v", err))
			return err
		}
		if len(lista) == 0 {
			svc.Publisher().SendTextMessage(chatID, "Nenhum rastreamento adicionado")
			return nil
		}
		remover := make([]string, len(lista))
		for i, rastreio := range lista {
			remover[i] = helpers.BotMenuOption{Command: "correio",
				Caption: rastreio.String(),
				Value:   fmt.Sprintf("remover:%s", rastreio.CodObjeto),
			}.String()
		}

		svc.Publisher().SendMenuKeyboard(update.CallbackQuery.Message.Chat.ID, "Remover rastreamento", "correio", remover...)
		return nil
	}

	if ok, codObjeto, err := runRemoverRastreamento(chatID, update.CallbackQuery.Data); ok {
		if err != nil {
			svc.Publisher().SendTextMessage(chatID, fmt.Sprintf("%v", err))
		} else {
			svc.Publisher().SendTextMessage(chatID, fmt.Sprintf("Rastreamento removido - %s", codObjeto))
		}
		return nil
	}

	return nil
}

func runRemoverRastreamento(chatID int64, comando string) (success bool, codObjeto string, err error) {
	if !strings.HasPrefix(comando, "correio|remover:") {
		return false, "", errors.New("comando inválido")
	}
	success = true
	codObjeto = strings.TrimPrefix(comando, "correio|remover:")
	err = correiosService.repository.RemoveRastreamento(chatID, codObjeto)
	if err != nil {
		err = fmt.Errorf("Erro ao remover rastreamento [%s] %v", codObjeto, err)
	}
	return
}

func correiosRastrearEncomenda(ctx context.Context, u tgbotapi.Update) (err error) {
	svc := tgbot.GetBotService(ctx)
	codRastreamento := strings.Trim(u.Message.CommandArguments(), " ")
	if len(codRastreamento) == 0 {
		svc.Publisher().SendMenuKeyboard(u.Message.Chat.ID, "Correios", "correio", "Listar rastreamentos:listar", "-", "Remover rastreamentos:remover")
		return
	}
	estadoAnterior := correiosService.repository.GetRastreio(codRastreamento, u.Message.Chat.ID)
	if estadoAnterior != nil && estadoAnterior.CodObjeto == codRastreamento {
		svc.Publisher().ReplyToMessage(u, "Rastreamento já foi adicionado")
		return
	}
	_, ultimoEstado, err := correiosService.VerificarRastreamento(codRastreamento, u.Message.Chat.ID)
	if err != nil {
		svc.Publisher().ReplyToMessage(u, fmt.Sprintf("%v", err))
	} else {
		svc.Publisher().ReplyToMessage(u, fmt.Sprintf("Rastreamento %s", ultimoEstado.String()))
	}
	return nil
}

func correiosVerificarRastreios(ctx context.Context) (err error) {
	svc := tgbot.GetBotService(ctx)

	rastreios, err := correiosService.repository.GetRastreiosPendentes()
	if err == nil && len(rastreios) == 0 {
		err = fmt.Errorf("Nenhum rastreamento pendente")
	}
	if err != nil {
		return
	}

	for _, rastreio := range rastreios {
		atualizado, err := correiosVerificarRastreio(rastreio)
		if err == nil && atualizado {
			svc.Publisher().SendTextMessage(rastreio.ChatID, fmt.Sprintf("Rastreamento %s", rastreio.String()))
		}
	}
	return
}

func correiosVerificarRastreio(rastreio *CorreiosRastreioModel) (atualizado bool, err error) {
	atualizado, ultimoEstado, err := correiosService.VerificarRastreamento(rastreio.CodObjeto, rastreio.ChatID)
	if atualizado {
		rastreio.DataHoraEvento = ultimoEstado.DataHoraEvento
		rastreio.DescEvento = ultimoEstado.DescEvento
		rastreio.ObjetoEntregue = ultimoEstado.ObjetoEntregue
	}
	return atualizado, err
}
