test:
	@echo "+ $@"
	go test -v -timeout 1m -failfast -race -coverprofile=profile.txt -covermode=atomic -count=1 ./...

tidy:
	@echo "+ $@"
	go mod tidy
