package menu

import (
	"docker_tray/docker"
	"docker_tray/icon"
	"docker_tray/logger"
	"fmt"
	"strings"
	"time"

	"github.com/getlantern/systray"
)

type MenuItem struct {
	Container *docker.Container
	Button    *systray.MenuItem
	CloseChan chan struct{}
}

type Menu struct {
	logger.Logger
	dockerService *docker.DockerService
	MenuItems     []MenuItem
	QuitItem      *systray.MenuItem
}

func NewMenu(dockerService *docker.DockerService) *Menu {
	target := Menu{}

	target.dockerService = dockerService
	return &target
}

func (m *Menu) Init() []MenuItem {
	systray.SetTemplateIcon(icon.Data, icon.Data)

	m.QuitItem = systray.AddMenuItem("<[Quit]>", "Quit the whole app")
	systray.AddSeparator()

	containers, err := m.dockerService.GetAllContainerList()
	if err != nil {
		m.LogError("Ошибка получения списка контейнеров", err)
		systray.Quit()
	}

	for i, c := range containers {
		m.LogInfo(fmt.Sprintf("Menu Init: %10s %20s %10s", c.ID[:10], c.Image, c.State))
		var mEnable *systray.MenuItem

		if strings.Compare(c.State, "running") == 0 {
			mEnable = systray.AddMenuItemCheckbox(c.Image, c.Image, true)
		} else {
			mEnable = systray.AddMenuItemCheckbox(c.Image, c.Image, false)
		}

		m.MenuItems = append(m.MenuItems, MenuItem{
			Container: &containers[i],
			Button:    mEnable,
			CloseChan: make(chan struct{}),
		})
	}

	for i, v := range m.MenuItems {
		m.LogInfo(fmt.Sprintf("Item %d: %s", i, v.Container.Image))
	}
	return m.MenuItems
}

func (m *Menu) Start(done chan struct{}) {
	// Обработка нажатий на пункты меню
	for _, item := range m.MenuItems {
		m.startMenuItem(done, item)
	}

	// Обработка нажатия пункт Quit
	go func() {
		<-m.QuitItem.ClickedCh
		m.LogInfo("Quit")
		systray.Quit()
	}()

	// Переодическое обновление списка меню
	m.startRefresherMenu(done)
}

func (m *Menu) addMenuItem(c docker.Container) MenuItem {

	var mEnable *systray.MenuItem

	if strings.Compare(c.State, "running") == 0 {
		mEnable = systray.AddMenuItemCheckbox(c.Image, c.Image, true)
	} else {
		mEnable = systray.AddMenuItemCheckbox(c.Image, c.Image, false)
	}

	res := MenuItem{
		Container: &c,
		Button:    mEnable,
		CloseChan: make(chan struct{}),
	}
	m.MenuItems = append(m.MenuItems, res)
	return res
}

func (m *Menu) closeMenuItem(i int) {
	m.MenuItems[i].Button.Hide()
	close(m.MenuItems[i].CloseChan)
	m.MenuItems = append(m.MenuItems[:i], m.MenuItems[i+1:]...)
}

func (m *Menu) startMenuItem(done chan struct{}, item MenuItem) {

	go func() {
		var err error

		for {
			select {
			case <-item.Button.ClickedCh:
				if item.Container.State == "running" {
					err = m.dockerService.ContainerStop(item.Container.ID)
					if err != nil {
						m.LogError("Ошибка остановки контейнера: "+item.Container.Image, err)
					} else {
						m.LogInfo("Остановлен контейнер: " + item.Container.Image)
						item.Button.Uncheck()
					}
				} else {
					err = m.dockerService.ContainerStart(item.Container.ID)
					if err != nil {
						m.LogError("Ошибка запуска контейнера: "+item.Container.Image, err)
					} else {
						m.LogInfo("Запущен контейнер: " + item.Container.Image)
						item.Button.Check()
					}
				}
			case <-item.CloseChan:
				m.LogInfo("Остановлена горутина для: " + item.Container.Image)
				return
			case <-done:
				m.LogInfo("Закрыта горутина для: " + item.Container.Image)
				return
			}
		}
	}()
}

func (m *Menu) startRefresherMenu(done chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				m.LogInfo("Закрыта горутина обновления меню")
				return
			case <-time.After(time.Second * 2):
				containerMap, err := m.dockerService.GetAllContainerMap()
				if err != nil {
					m.LogInfo("Закрыта горутина обновления меню")
					break
				}

				for i := 0; i < len(m.MenuItems); i++ {
					c, ok := containerMap[m.MenuItems[i].Container.ID]
					if !ok {
						m.closeMenuItem(i)
						i--
						continue
					}
					delete(containerMap, m.MenuItems[i].Container.ID)
					m.MenuItems[i].Container.State = c.State

					if m.MenuItems[i].Container.State == "running" {
						m.MenuItems[i].Button.Check()
					} else {
						m.MenuItems[i].Button.Uncheck()
					}
				}

				if len(containerMap) > 0 {
					for _, value := range containerMap {
						m.LogInfo("Добавился контейнер: " + value.Image)
						newItem := m.addMenuItem(value)
						m.startMenuItem(done, newItem)
					}
				}
			}
		}
	}()
}
