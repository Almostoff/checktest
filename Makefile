start-docker:
		docker-compose up

stop-docker:
		docker-compose stop

#test:
#		cd server_sub && go test ./cache && go test ./database

run-subscriber:
		cd server_sub && go build ./
		cd server_sub && server_sub.exe

run-publisher:
		cd orderPub && go build ./
		cd orderPub && orderPub.exe

.PHONY: start-docker stop-docker run-subscriber run-publisher