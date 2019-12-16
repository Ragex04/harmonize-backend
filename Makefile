all:
	-mkdir bin
	go build -o bin/server main.go

run:
	./bin/server

client:
	cd ../harmonize-frontend; yarn build
	-rm -rf ./wwwroot/*
	touch ./wwwroot/.gitkeep
	cp -r ../harmonize-frontend/build/* ./wwwroot