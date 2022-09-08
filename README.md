# docker_tray
Tray indicator for start/stop docker containers.
Script start docker services:
- docker
- docker.socket
- containerd

And stop it with quit.
`Updating list of containers on show`

## Dependencies:window
- [gotk3](https://github.com/gotk3/gotk3/gtk)
- [go-appindicator](https://github.com/dawidd6/go-appindicator)

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

## Fix
if not found `appindicator`
```
./appindicator_fix.sh
```
make symlinks to `ayatana-appindicator`
