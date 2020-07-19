GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
COUNT=$$(docker ps -a -q)

fmt:
	gofmt -w $(GOFMT_FILES)

build: fmt
	docker build -f docker/Dockerfile -t ssastry22/atlas-repo:photo .
	docker build -f psql/Dockerfile -t ssastry22/atlas-repo:psql .
	docker build -f cmd/consumer/Dockerfile -t ssastry22/atlas-repo:consumer .
	cd docker; docker-compose up -d

clean:
	cd docker; docker-compose down
	docker volume rm -f $$(docker volume ls  -q)
	docker image rm -f $$(docker image ls -a -q)