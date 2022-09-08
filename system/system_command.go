package system

import (
	"os/exec"
	"strings"

	"docker_tray/logger"
)

var pass string
var l logger.Logger

func GetPass() string {
	if pass == "" {
		l.LogInfo("Получение пароля")
		out, _ := exec.Command("/bin/bash", "get_pass.sh").Output()
		pass = strings.TrimSpace(string(out))
	}

	return pass
}
