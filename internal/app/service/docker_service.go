package service

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"docker-tray/internal/app/logger"
	"docker-tray/internal/app/obj"
)

type DockerService struct {
	logger.Logger
	ctx                    context.Context
	client                 *client.Client
	passService            PassService
	isControlDockerService bool
}

func NewDockerService(ctx context.Context, passService *PassService, isControlDockerService bool) (*DockerService, error) {
	var err error

	target := DockerService{}
	target.passService = *passService
	target.ctx = ctx
	target.isControlDockerService = isControlDockerService

	status, err := target.checkDockerService()
	if err != nil {
		target.LogError("Ошибка проверки статуса Docker", err)
		return &target, err
	}
	if status == "inactive" {
		if target.isControlDockerService {
			err = target.startDockerService()
			if err != nil {
				target.LogError("Ошибка при запуске docker ", err)
				return nil, err
			}
		} else {
			return &target, errors.New("Не запущен Docker")
		}
	}

	target.client, err = client.NewClientWithOpts()
	if err != nil {
		target.LogError("Ошибка создания клиента docker", err)
		return nil, err
	}
	return &target, err
}

func (d *DockerService) GetAllContainerList() ([]obj.Container, error) {

	containers, err := d.client.ContainerList(d.ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		d.LogError("Ошибка получения списка контейнеров", err)
		return nil, err
	}

	var res []obj.Container
	for _, container := range containers {
		// d.LogInfo(fmt.Sprintf("%10s %20s %10s", container.ID[:10], container.Image, container.State))
		res = append(res, obj.Container{
			ID:    container.ID,
			Image: container.Image,
			State: strings.TrimSpace(container.State),
		})
	}
	return res, nil
}

func (d *DockerService) getActiveContainerList() ([]obj.Container, error) {

	containers, err := d.client.ContainerList(d.ctx, types.ContainerListOptions{})
	if err != nil {
		d.LogError("Ошибка получения списка контейнеров", err)
		return nil, err
	}

	var res []obj.Container
	for _, container := range containers {
		// d.LogInfo(fmt.Sprintf("%10s %20s %10s", container.ID[:10], container.Image, container.State))
		res = append(res, obj.Container{
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

func (d *DockerService) stopActiveContainers() {
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
}

func (d *DockerService) Close() {
	var err error

	d.stopActiveContainers()
	d.client.Close()

	if d.isControlDockerService {
		err = d.stopDockerService()
		if err != nil {
			d.LogError("Ошибка при остановке docker", err)
		}
	}
}

func (d *DockerService) checkDockerService() (string, error) {
	var err error
	var out []byte

	out, err = exec.Command("systemctl", "is-active", "docker").Output()
	outString := strings.TrimSpace(string(out))
	if err != nil {
		if outString != "inactive" {
			d.LogError("Ошибка проверки is-active docker "+string(out), err)
			return "", err
		}
		d.LogInfo(fmt.Sprintf("Docker не запущен %s %s", string(out), err.Error()))
		err = nil
	}

	return outString, err
}

func (d *DockerService) startDockerService() error {
	var err error
	var out []byte
	var cmd *exec.Cmd

	password, err := d.passService.GetPass()
	if err != nil {
		d.LogInfo("Отмена ввода пароля")
		return errors.New("Отмена ввода пароля")
	}
	cmd = exec.Command("sudo", "-S", "systemctl", "start", "docker", "docker.socket", "containerd")
	cmd.Stdin = strings.NewReader(password)

	out, err = cmd.Output()
	if err != nil {
		d.LogError("Ошибка запуска docker "+string(out), err)
		return err
	}
	d.LogInfo("Запущены сервисы docker")

	return nil
}

func (d *DockerService) stopDockerService() error {
	password, err := d.passService.GetPass()
	if err != nil {
		d.LogError("Отмена ввода пароля", err)
		return errors.New("Отмена ввода пароля")
	}

	cmd := exec.Command("sudo", "-S", "systemctl", "stop", "docker", "docker.socket", "containerd")
	cmd.Stdin = strings.NewReader(password)
	out, err := cmd.Output()
	if err != nil {
		d.LogError("Ошибка остановки docker: "+string(out), err)
		return err
	}

	d.LogInfo("Остановлены сервисы docker")

	return nil
}
