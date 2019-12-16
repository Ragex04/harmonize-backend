all:
	-mkdir bin
	go build -race -o bin/server main.go

run:
	./bin/server

client:
	cd ../harmonize-frontend; yarn build
	cp -r ../harmonize-frontend/build/* ./wwwroot