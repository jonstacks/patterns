test:
	go test -v -race -coverprofile=coverage.txt -timeout=30s -covermode=atomic ./...