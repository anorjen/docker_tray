package main

import (
	"context"
	"docker_tray/docker"
	"docker_tray/logger"

	appindicator "github.com/dawidd6/go-appindicator"
	"github.com/gotk3/gotk3/gtk"
)

var dockerService *docker.DockerService
var log logger.Logger

func close() {
	gtk.MainQuit()

	if dockerService != nil {
		dockerService.Close()
	}
}

func main() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dockerService, err = docker.NewDockerService(ctx)
	if err != nil {
		log.LogError("Ошибка при создании docker service", err)
		close()
		return
	}

	gtk.Init(nil)
	createTrayMenu()
	gtk.Main()
}

func createTrayMenu() {

	menu, err := gtk.MenuNew()
	if err != nil {
		log.LogError("Ошибка создания меню", err)
	}

	indicator := appindicator.New("indicator", "./resources/docker_icon.svg", appindicator.CategoryApplicationStatus)
	indicator.SetTitle("")
	indicator.SetLabel("docker_tray", "")
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	var items []*gtk.MenuItem
	menu.Connect("show", func() {
		items = items[:0]
		var i uint
		currentItems := menu.GetChildren()
		for i = 0; i < currentItems.Length(); i++ {
			w := currentItems.NthData(i).(*gtk.Widget)
			menu.Remove(w)
		}

		containers, _ := dockerService.GetAllContainerList()
		for _, c := range containers {
			item, err := gtk.MenuItemNewWithLabel(c.Image)
			if err != nil {
				log.LogError("Ошибка создания элемента меню", err)
			}
			if c.State == "running" {
				item.SetOpacity(1)
			} else {
				item.SetOpacity(0.5)
			}
			items = append(items, item)
		}
		for i, _ := range items {

			menu.Add(items[i])
			j := i
			items[i].Connect("activate", func() {

				if containers[j].State == "running" {
					err = dockerService.ContainerStop(containers[j].ID)
					if err != nil {
						log.LogError("Ошибка остановки контейнера: "+containers[j].Image, err)
					} else {
						log.LogInfo("Остановлен контейнер: " + containers[j].Image)
						items[j].SetOpacity(0.5)
					}
				} else {
					err = dockerService.ContainerStart(containers[j].ID)
					if err != nil {
						log.LogError("Ошибка запуска контейнера: "+containers[j].Image, err)
					} else {
						log.LogInfo("Запущен контейнер: " + containers[j].Image)
						items[j].SetOpacity(1)
					}
				}
			})
		}

		item, err := gtk.MenuItemNewWithLabel("Quit")
		if err != nil {
			log.LogError("Ошибка создания элемента меню Quit", err)
		}
		item.Connect("activate", func() {
			close()
		})

		item.SetMarginTop(10)
		menu.Add(item)

		menu.ShowAll()
	})
}
