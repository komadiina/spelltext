package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gdamore/tcell/v2"
	"github.com/komadiina/spelltext/client/audio"
	"github.com/komadiina/spelltext/client/config"
	"github.com/komadiina/spelltext/client/constants"
	"github.com/komadiina/spelltext/client/factory"
	"github.com/komadiina/spelltext/client/hooks"
	"github.com/komadiina/spelltext/client/types"
	"github.com/komadiina/spelltext/shared"
	"github.com/komadiina/spelltext/utils/singleton/logging"
	"github.com/rivo/tview"

	pbRepo "github.com/komadiina/spelltext/proto/repo"
)

func init() {}

func main() {
	var level log.Level = log.DebugLevel
	logging.Init(level, "client", true)
	logger := logging.Get("client", true)

	err := os.Setenv("CONFIG_FILE", "config.yml")
	logger.Debug("loading client config...", "CONFIG_FILE", os.Getenv("CONFIG_FILE"))
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("failed to load config, using default values.", "reason", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// init audio manager
	am := audio.NewManager(44000)
	err = am.Preload(constants.PRELOAD)

	if err != nil {
		logger.Fatal("failed to preload audio files", "reason", err)
	}

	am.AudioEnabled = cfg.AudioEnabled

	client := types.SpelltextClient{
		Config:       cfg,
		Logger:       logger,
		App:          tview.NewApplication(),
		Context:      &ctx,
		AppStorage:   make(map[string]any),
		AudioManager: am,
		Storage: &types.AppStorage{
			Ministate:         &types.Ministate{},
			CurrentUser:       &pbRepo.User{},
			SelectedCharacter: &pbRepo.Character{},
			SelectedVendor:    &pbRepo.Vendor{},
			EquipSlots:        nil,
		},
	}

	client.AudioManager.PlayBackground(logger)

	logger.Debug("initializing clients...")
	hooks.InitializeClients(&client)
	defer hooks.CloseClients(&client)
	logger.Debug("clients initialized.")

	logger.Debug("initializng nats...")
	nc, _, err := hooks.InitializeNats(cfg)
	if err != nil {
		logger.Fatal("failed to init nats/jetstream", "reason", err)
	}

	client.Nats = nc
	logger.Debug("nats/js initialized.")

	logger.Debug("initializing PageManager factory...")
	client.PageManager = factory.NewPageManager(client.Logger, client.App)
	hooks.InitializePages(&client)
	logger.Debug("PageManager factory initialized.")

	client.NavigateTo = func(pageKey string) {
		if client.PageManager.HasPage(pageKey) == false {
			logger.Fatal("page not found", "page", pageKey)
			return
		}

		client.PageManager.Push(pageKey, false)
	}

	client.NavigateTo(constants.PAGE_LOGIN)

	client.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			client.AudioManager.Play(constants.BLIP_BACKWARD, client.Logger)

			if client.PageManager.Pop() == -1 {
				client.App.Stop()
				return nil
			}
		} else if event.Key() == tcell.KeyEnter {
			client.AudioManager.Play(constants.BLIP_FORWARD, client.Logger)
		} else {
			client.AudioManager.Play(constants.BLIP_INPUT, client.Logger)
		}

		return event
	})

	if err := client.App.SetRoot(client.PageManager.Pages, true).EnableMouse(true).Run(); err != nil {
		client.Logger.Error(err)
	}

	// cleanup
	client.Nats.Drain()
	defer cancel()
	logger.Debug("client shutdown.")

	fmt.Print(shared.BANNER)
	goodbye := `
> thanks for playing this torturefest
> a game by ogg/komadiina (https://github.com/komadiina)
> follow the development at https://github.com/komadiina/spelltext
~ kthxb
`
	fmt.Print(goodbye)
}
