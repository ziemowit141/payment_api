install_swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger:
	@echo Ensure you have the swagger CLI or this command will fail.
	@echo You can install the swagger CLI with: go get -u github.com/go-swagger/go-swagger/cmd/swagger
	@echo ....

	swagger generate spec -o ./swagger.yaml --scan-models

postgres_docker:
	@echo Start postgres database container
	docker run --name some-postgres  -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=myapp -e POSTGRES_SERVER=postgres -d postgres

build_server:
	@echo Build payment_api server
	go mod tidy
	go build -o server .

quicksetup:
	make postgres_docker

	@echo Wait 10 seconds for postgres docker to start
	sleep 10

	make build_server
	./server

tests:
	@echo Tests won't complete without database setup
	ginkgo run -r .