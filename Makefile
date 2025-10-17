run:
	go run cmd/server/main.go
test:
	go test -v ./...
cover:
	go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out
bench:
	go test -bench . ./router -benchmem -benchtime=1000x
