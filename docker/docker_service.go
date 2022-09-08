package docker

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"docker_tray/logger"
	"docker_tray/system"
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
	var out []byte
	var cmd *exec.Cmd
	target := DockerService{}
	target.ctx = ctx

	out, err = exec.Command("systemctl", "is-active", "docker").Output()
	outString := strings.TrimSpace(string(out))
	if err != nil {
		if outString == "inactive" {
			target.LogInfo(fmt.Sprintf("Docker не запущен %s %s", string(out), err.Error()))
			err = nil
		} else {
			target.LogError("Ошибка проверки is-active docker "+string(out), err)
			return &target, err
		}
	}
	if outString == "inactive" {
		cmd = exec.Command("sudo", "-S", "systemctl", "start", "docker", "docker.socket", "containerd")
		cmd.Stdin = strings.NewReader(system.GetPass())

		out, err = cmd.Output()
		if err != nil {
			target.LogError("Ошибка запуска docker "+string(out), err)
			return &target, err
		}
		target.LogInfo("Запущены сервисы docker")
	}

	target.client, err = client.NewClientWithOpts()
	if err != nil {
		target.LogError("Ошибка создания клиента docker", err)
		return nil, err
	}
	return &target, err
}

func (s *DockerService) GetAllContainerList() ([]Container, error) {

	containers, err := s.client.ContainerList(s.ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		s.LogError("Ошибка получения списка контейнеров", err)
		return nil, err
	}

	var res []Container
	for _, container := range containers {
		s.LogInfo(fmt.Sprintf("%10s %20s %10s", container.ID[:10], container.Image, container.State))
		res = append(res, Container{
			ID:    container.ID,
			Image: container.Image,
			State: strings.TrimSpace(container.State),
		})
	}
	return res, nil
}

func (s *DockerService) GetAllContainerMap() (map[string]Container, error) {
	containers, err := s.client.ContainerList(s.ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		s.LogError("Ошибка получения списка контейнеров", err)
		return nil, err
	}

	var containerMap = make(map[string]Container)
	for _, c := range containers {
		containerMap[c.ID] = Container{
			ID:    c.ID,
			Image: c.Image,
			State: strings.TrimSpace(c.State),
		}
	}

	return containerMap, nil
}

func (s *DockerService) getActiveContainerList() ([]Container, error) {

	containers, err := s.client.ContainerList(s.ctx, types.ContainerListOptions{})
	if err != nil {
		s.LogError("Ошибка получения списка контейнеров", err)
		return nil, err
	}

	var res []Container
	for _, container := range containers {
		s.LogInfo(fmt.Sprintf("%10s %20s %10s", container.ID[:10], container.Image, container.State))
		res = append(res, Container{
			ID:    container.ID,
			Image: container.Image,
			State: strings.TrimSpace(container.State),
		})
	}
	return res, nil
}

func (s *DockerService) ContainerStart(containerID string) error {
	return s.client.ContainerStart(s.ctx, containerID, types.ContainerStartOptions{})
}

func (s *DockerService) ContainerStop(containerID string) error {
	return s.client.ContainerStop(s.ctx, containerID, nil)
}

func (s *DockerService) Close() {
	var err error

	containers, err := s.getActiveContainerList()
	if err != nil {
		s.LogError("Ошибка получения активных контейнеров", err)
	} else {
		for _, c := range containers {
			if err = s.ContainerStop(c.ID); err != nil {
				s.LogError("Ошибка остановки контейнера "+c.Image, err)
			} else {
				s.LogInfo("Остановлен контейнер " + c.Image)
			}
		}
	}
	s.client.Close()

	cmd := exec.Command("sudo", "-S", "systemctl", "stop", "docker", "docker.socket", "containerd")
	cmd.Stdin = strings.NewReader(system.GetPass())
	out, err := cmd.Output()
	if err != nil {
		s.LogError("Ошибка остановки docker: "+string(out), err)
	} else {
		s.LogInfo("Остановлены сервисы docker")

	}
}
