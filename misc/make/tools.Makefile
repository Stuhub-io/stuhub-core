install-air:
	go install github.com/air-verse/air@latest

install-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	
install-golang-migrate:	
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

install-gorm-gentool:
	go install gorm.io/gen/tools/gentool@latest

.PHONY: install-golangci-lint install-air install-golang-migrate install-gorm-gentool
