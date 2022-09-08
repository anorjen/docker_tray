#!/bin/bash

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

echo $(REQUEST_PASS)
