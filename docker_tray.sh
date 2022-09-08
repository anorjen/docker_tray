#!/bin/bash

## Disable docker services
# sudo systemctl stop docker docker.socket containerd
# sudo systemctl disable docker docker.socket containerd

function REQUEST_PASS() {
	local var=1
	
	while [ $var -gt 0 ]
	do
		## Request password
		PASS=$(yad --entry --entry-text="pass" --hide-text)
		#~ echo $PASS
		
		## make sure to ask for password on next sudo
		sudo -k
		
		## check password
		if sudo -lS &>/dev/null << EOF
$PASS
EOF
		then
			var=0
		else
			var=1
		fi
	done
	
	echo $PASS
}

function START_MENU() {

	PASS=$(REQUEST_PASS)
	
	## Start docker
    if [[ "$(systemctl is-active docker)" == "inactive" ]]; then
		echo $PASS | sudo -S systemctl start docker docker.socket containerd
	fi
	
	PIPE_FIFO=$(mktemp -u /tmp/menutray.XXXXXXXX)
	
	function CREATE_PIPE_FIFO() {
		## 1 Create PIPE_FIFO file
		mkfifo $PIPE_FIFO

		## 2 Attach a filedescriptor to this PIPE_FIFO
		exec 3<> $PIPE_FIFO
	}

	function SET_MENU_LIST(){
		local GET_MENU="[Quit]! bash -c QUIT $PASS"
		local LIST=(`docker container ls -a --format "{{.Names}}"`)
		
		for i in "${LIST[@]}"
		do
			GET_MENU="[stop] $i! docker stop $i|$GET_MENU";
			GET_MENU="[start] $i! docker start $i|$GET_MENU";
		done
		
		echo $GET_MENU
	}
	
	## Action on left mouse click
	function QUIT() {
		exec 3<> $PIPE_FIFO
		echo "quit" >&3
		rm -f $PIPE_FIFO
		
		## Stop docker
		docker stop $(docker container ls -a -q)
		echo $1 | sudo -S systemctl stop docker docker.socket containerd
	}
	export -f QUIT
	export PIPE_FIFO
	
	MENU_ITEMS=$(SET_MENU_LIST)
	
	## Defaults
	TRAY_ICON="$HOME/docker_tray/docker_icon.png"
	POPUP_TEXT="docker tray"
	
	CREATE_PIPE_FIFO
	
	## 3 Run yad and tell it to read its stdin from the file descriptor
	GUI=$(yad --notification --kill-parent --listen \
	--image="$TRAY_ICON" \
	--text="$POPUP_TEXT" \
	--command=<&3) & 

	## 4 Write menu to file descriptor to generate MENU
	echo "menu:$MENU_ITEMS" >&3
}

## START
START_MENU
