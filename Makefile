BUILDPATH=$(CURDIR)
NAME=golang-jobs-scrapper

build: 
	@echo "Creando Binario ..."
	@go build -mod=vendor -ldflags '-s -w' -o $(BUILDPATH)/build/bin/${NAME} cmd/${NAME}/main.go
	@echo "Binario generado en build/bin/${NAME}"

test: 
	@echo "Ejecutando tests..."
	@go test ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out

coverage: 
	@echo "Coverfile..."
	@go test ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out
	@go tool cover -html=coverfile_out -o coverfile_out.html

docker:
	@docker build . -t $(NAME):latest -f iaas/Dockerfile

vendor:
	@echo "Vendoring..."
	@go mod vendor
