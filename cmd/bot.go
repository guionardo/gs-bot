package cmd

import (
	"github.com/guionardo/go-tgbot/tgbot"
	"github.com/guionardo/go-tgbot/tgbot/automations"
	"github.com/guionardo/gs-bot/services/correios"
	"github.com/guionardo/gs-bot/services/links"
)

func GetBot() *tgbot.GoTGBotService {
	svc := tgbot.CreateBotService().
		LoadConfigurationFromEnv("TG_").
		InitBot()

	automations.AddStartupGreetingsAutomation(svc)
	automations.AddSetupCommandsAutomation(svc)
	automations.AddHouseKeepingAutomation(svc)

	correios.SetupBot(svc)
	links.SetupBot(svc)
	
	return svc
}
