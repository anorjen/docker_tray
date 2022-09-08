package main

import (
	"context"

	"github.com/getlantern/systray"

	"docker_tray/docker"
	"docker_tray/logger"
	"docker_tray/menu"
)

func main() {
	l := logger.Logger{}
	ctx := context.Background()
	done := make(chan struct{})

	dockerService, err := docker.NewDockerService(ctx)
	if err != nil {
		l.LogError("Ошибка создания dockerService", err)
		return
	}

	menu := menu.NewMenu(dockerService)

	onExit := func() {
		close(done)
		dockerService.Close()
		l.LogInfo("Завершение работы")
	}

	onReady := func() {
		menu.Init()
		menu.Start(done)
	}

	systray.Run(onReady, onExit)

}
