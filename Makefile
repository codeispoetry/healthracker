deploy-gui:
	rsync -avz --delete index.html manifest.json sw.js icon.svg icon-512.svg favicon.svg tom-rose.de:./httpdocs/healthtracker

deploy-server:
	go build -o healthtracker .
	rsync -avz --delete healthtracker tom-rose.de:./healthtracker
