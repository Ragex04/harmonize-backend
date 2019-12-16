all:
	-mkdir bin
	go build -race -o bin/server main.go