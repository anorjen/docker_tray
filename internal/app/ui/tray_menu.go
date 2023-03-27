package ui

import (
	"docker-tray/internal/app/logger"
	"docker-tray/internal/app/obj"

	"github.com/dawidd6/go-appindicator"
	"github.com/gotk3/gotk3/gtk"
)

const (
	iconFile = "./assets/docker_icon.svg"
)

type DockerService interface {
	ContainerStart(containerID string) error
	ContainerStop(containerID string) error
	GetAllContainerList() ([]obj.Container, error)
	Close()
}

type TrayMenu struct {
	logger.Logger
	dockerService DockerService
}

func NewTrayMenu(dockerService DockerService) *TrayMenu {
	return &TrayMenu{dockerService: dockerService}
}

func (t *TrayMenu) CreateTrayMenu() {
	gtk.Init(nil)

	menu, err := gtk.MenuNew()
	if err != nil {
		t.LogError("Ошибка создания меню", err)
	}

	indicator := appindicator.New("indicator", iconFile, appindicator.CategoryApplicationStatus)
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

		containers, _ := t.dockerService.GetAllContainerList()
		for _, c := range containers {
			item, err := gtk.MenuItemNewWithLabel(c.Image)
			if err != nil {
				t.LogError("Ошибка создания элемента меню", err)
			}
			if c.State == "running" {
				item.SetOpacity(1)
			} else {
				item.SetOpacity(0.5)
			}
			items = append(items, item)
		}

		for i := range items {
			menu.Add(items[i])
			j := i
			items[i].Connect("activate", func() {

				if containers[j].State == "running" {
					err = t.dockerService.ContainerStop(containers[j].ID)
					if err != nil {
						t.LogError("Ошибка остановки контейнера: "+containers[j].Image, err)
					} else {
						t.LogInfo("Остановлен контейнер: " + containers[j].Image)
						items[j].SetOpacity(0.5)
					}
				} else {
					err = t.dockerService.ContainerStart(containers[j].ID)
					if err != nil {
						t.LogError("Ошибка запуска контейнера: "+containers[j].Image, err)
					} else {
						t.LogInfo("Запущен контейнер: " + containers[j].Image)
						items[j].SetOpacity(1)
					}
				}
			})
		}

		item, err := gtk.MenuItemNewWithLabel("Quit")
		if err != nil {
			t.LogError("Ошибка создания элемента меню Quit", err)
		}
		item.Connect("activate", func() {
			t.close()
		})

		item.SetMarginTop(10)
		menu.Add(item)

		menu.ShowAll()
	})

	gtk.Main()
}

func (t *TrayMenu) close() {
	gtk.MainQuit()

	if t.dockerService != nil {
		t.dockerService.Close()
	}
}