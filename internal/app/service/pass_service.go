package service

import (
	"docker-tray/internal/app/logger"
	"docker-tray/internal/app/ui"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type PassService struct {
	logger.Logger
	password string
}

func NewPassService() *PassService {
	return &PassService{}
}

func (s *PassService) GetPass() (string, error) {
	if s.password == "" {
		passWindow := ui.NewPassWindow()
		// passWindow.Close()

		var pass string
		var cancel bool
		for {
			pass, cancel = passWindow.GetPass()
			if !cancel {
				if ok := s.checkPassword(pass); !ok {
					continue
				}
			}

			s.password = pass
			break
		}

		// passWindow.Close()

		if cancel {
			s.LogInfo("Отмена ввода пароля")
			return "", errors.New("Отмена ввода пароля")
		}
	}

	return s.password, nil
}

func (s *PassService) checkPassword(password string) bool {
	cmd := exec.Command("sudo", "-lkS")
	cmd.Stdin = strings.NewReader(password)
	err := cmd.Run()
	if err != nil {
		s.LogError("Неверный пароль", err)
	}
	code := cmd.ProcessState.ExitCode()
	s.LogInfo(fmt.Sprintf("Check password code: %d", code))

	return code == 0
}
