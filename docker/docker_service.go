package docker

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"docker_tray/logger"
	"docker_tray/window"
)

type Container struct {
	ID    string
	Image string
	State string
}

type DockerService struct {
	logger.Logger
	client *client.Client
	ctx    context.Context
}

func NewDockerService(ctx context.Context) (*DockerService, error) {
	var err error
	target := DockerService{}
	target.ctx = ctx

	err = target.startDockerService()
	if err != nil {
		target.LogError("Ошибка при запуске docker ", err)
		return nil, err
	}

	target.client, err = client.NewClientWithOpts()
	if err != nil {
		target.LogError("Ошибка создания клиента docker", err)
		return nil, err
	}
	return &target, err
}

func (d *DockerService) GetAllContainerList() ([]Container, error) {

	containers, err := d.client.ContainerList(d.ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		d.LogError("Ошибка получения списка контейнеров", err)
		return nil, err
	}

	var res []Container
	for _, container := range containers {
		// d.LogInfo(fmt.Sprintf("%10s %20s %10s", container.ID[:10], container.Image, container.State))
		res = append(res, Container{
			ID:    container.ID,
			Image: container.Image,
			State: strings.TrimSpace(container.State),
		})
	}
	return res, nil
}

func (d *DockerService) getActiveContainerList() ([]Container, error) {

	containers, err := d.client.ContainerList(d.ctx, types.ContainerListOptions{})
	if err != nil {
		d.LogError("Ошибка получения списка контейнеров", err)
		return nil, err
	}

	var res []Container
	for _, container := range containers {
		// d.LogInfo(fmt.Sprintf("%10s %20s %10s", container.ID[:10], container.Image, container.State))
		res = append(res, Container{
			ID:    container.ID,
			Image: container.Image,
			State: strings.TrimSpace(container.State),
		})
	}
	return res, nil
}

func (d *DockerService) ContainerStart(containerID string) error {
	return d.client.ContainerStart(d.ctx, containerID, types.ContainerStartOptions{})
}

func (d *DockerService) ContainerStop(containerID string) error {
	return d.client.ContainerStop(d.ctx, containerID, nil)
}

func (d *DockerService) Close() {
	var err error

	containers, err := d.getActiveContainerList()
	if err != nil {
		d.LogError("Ошибка получения активных контейнеров", err)
	} else {
		for _, c := range containers {
			if err = d.ContainerStop(c.ID); err != nil {
				d.LogError("Ошибка остановки контейнера "+c.Image, err)
			} else {
				d.LogInfo("Остановлен контейнер " + c.Image)
			}
		}
	}
	d.client.Close()

	err = d.stopDockerService()
	if err != nil {
		d.LogError("Ошибка при остановке docker", err)
	}
}

func (d *DockerService) startDockerService() error {
	var err error
	var out []byte
	var cmd *exec.Cmd

	out, err = exec.Command("systemctl", "is-active", "docker").Output()
	outString := strings.TrimSpace(string(out))
	if err != nil {
		if outString == "inactive" {
			d.LogInfo(fmt.Sprintf("Docker не запущен %s %s", string(out), err.Error()))
			err = nil
		} else {
			d.LogError("Ошибка проверки is-active docker "+string(out), err)
			return err
		}
	}
	if outString == "inactive" {

		password, cancel := window.GetPass()
		if !cancel {
			cmd = exec.Command("sudo", "-S", "systemctl", "start", "docker", "docker.socket", "containerd")
			cmd.Stdin = strings.NewReader(password)

			out, err = cmd.Output()
			if err != nil {
				d.LogError("Ошибка запуска docker "+string(out), err)
				return err
			}
			d.LogInfo("Запущены сервисы docker")
		} else {
			d.LogInfo("Отмена ввода пароля")
			return errors.New("Отмена ввода пароля")
		}
	}
	return nil
}

func (d *DockerService) stopDockerService() error {

	password, cancel := window.GetPass()
	if !cancel {
		cmd := exec.Command("sudo", "-S", "systemctl", "stop", "docker", "docker.socket", "containerd")
		cmd.Stdin = strings.NewReader(password)
		out, err := cmd.Output()
		if err != nil {
			d.LogError("Ошибка остановки docker: "+string(out), err)
			return err
		} else {
			d.LogInfo("Остановлены сервисы docker")
		}
	} else {
		d.LogInfo("Отмена ввода пароля")
		return errors.New("Отмена ввода пароля")
	}
	return nil
}
