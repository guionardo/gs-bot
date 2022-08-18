package links

import (
	"context"
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
	linksService    *LinksService
	linksRepository *LinksRepository
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
	linksRepository = CreateLinksRepository(db, infra.GetLogger("links_repo")).Init()
	linksService = &LinksService{
		repository: linksRepository,
		logger:     infra.GetLogger("links_service"),
	}

	// svc.AddCommandHandlers(
	// 	&tgbot.ListenerCommandHandler{
	// 		Command: "correio",
	// 		Title:   "Rastrear encomenda do Correio",
	// 		Func:    correiosRastrearEncomenda,
	// 	},
	// )
	// svc.AddCallbackHandlers(tgbot.CreateListenerCallbackHandler("Correios", "correio", correiosCallbackHandler))

	svc.AddHandlers(&tgbot.ListenerFilteredHandler{
		Title:  "link url",
		Filter: func(update tgbotapi.Update) bool { return helpers.IsValidUrl(strings.Trim(update.Message.Text, " ")) },
		Func:   registrarLink,
	})

	svc.SetupBackgroundSchedules(schedules.CreateSchedule("Verificar rastreamentos", time.Hour, linksNotificarNaoLidos))
}

func registrarLink(ctx context.Context, update tgbotapi.Update) error {
	link, err := linksService.FetchLink(update.Message.Text)
	if err == nil {
		link.ChatID = update.Message.Chat.ID
		err = linksService.repository.Save(link)
	}
	svc := tgbot.GetBotService(ctx)
	if err != nil {
		svc.Publisher().ReplyToMessage(update, "Erro ao registrar link: "+err.Error())
	} else {
		svc.Publisher().ReplyToMessage(update, "Link registrado com sucesso! "+link.Title)
	}
	return err
}

func linksNotificarNaoLidos(ctx context.Context) error {
	links, err := linksService.GetUnreaden()
	if err != nil {
		return err
	}
	svc := tgbot.GetBotService(ctx)
	for chatId, links := range links {
		msgLinks := make([]string, len(links))
		for i, link := range links {
			msgLinks[i] = link.Title + ":" + link.URL
		}
		svc.Publisher().SendMenuKeyboard(chatId, "Links não lidos", "link", msgLinks...)
	}
	return nil
}
