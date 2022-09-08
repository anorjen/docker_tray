# docker_tray
Tray indicator for start/stop docker containers.
Script start docker services:
- docker
- docker.socket
- containerd

And stop it with quit.
`Container list update by timing`

## Dependencies:
- [yad](https://sourceforge.net/projects/yad-dialog/) - request password window
- [systray](https://github.com/getlantern/systray) - for tray appindicator

## Build
```
go build
```

## Run
without `Build`
```
go run main.go
```
or after `Build`
```
./docker_tray
```
