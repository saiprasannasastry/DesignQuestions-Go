GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
COUNT=$$(docker ps -a -q)

help:
	@echo 'Usage: make <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    fmt                formats the source code.'
	@echo '    clean              gets the containers down and removes all docker images.'

fmt:

	go fmt ./...

build: fmt
	docker build -f docker/Dockerfile -t ssastry22/atlas-repo:huddl .
	docker build -f psql/Dockerfile -t ssastry22/atlas-repo:psql .
	cd docker; docker-compose up -d
clean:
	cd docker; docker-compose down
	docker volume rm -f $$(docker volume ls  -q)
	docker image rm -f $$(docker image ls -a -q)