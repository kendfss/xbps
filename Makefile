default: install

install:
	go install
	[ ! -f /usr/bin/xbps ] && sudo ln -s ~/go/bin/xbps /usr/bin/xbps
	
