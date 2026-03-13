resetdb:
	rm -f diary.db
	sqlite3 diary.db < schema.sql

deploy-gui:
	rsync -avz --delete index.html tom-rose.de:./httpdocs/healthtracker

deploy-server:
	go build -o healthtracker .
	rsync -avz --delete healthtracker schema.sql tom-rose.de:./healthtracker