package app

import (
	"context"
	"docker-tray/internal/app/logger"
	"docker-tray/internal/app/service"
	"docker-tray/internal/app/ui"
	"os"
	"strings"
)

var (
	log                    logger.Logger
	isControlDockerService bool
	dockerService          *service.DockerService
	passService            *service.PassService
)

func Main() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	isControlDockerService = true
	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "-d") {
		isControlDockerService = false
	}

	passService = service.NewPassService()
	dockerService, err = service.NewDockerService(ctx, passService, isControlDockerService)
	if err != nil {
		log.LogError("Ошибка при создании docker service", err)
		return
	}
	trayMenu := ui.NewTrayMenu(dockerService)

	trayMenu.CreateTrayMenu()
}
