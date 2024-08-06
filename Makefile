build:
	sudo docker-compose build
build-nc:
	docker-compose build --no-cache
run-logs:
	sudo docker-compose up 
run:
	sudo docker-compose up -d  
stop:
	sudo docker-compose stop
down:
	sudo docker-compose down
clear-vm:
	docker-compose down --volumes --remove-orphans

