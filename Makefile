all: build

build:
	go build -o image_server -ldflags "-X haoyu.love/ImageServer/version.Version=`git describe --tags --dirty=-dev` -X haoyu.love/ImageServer/version.Build=`git rev-parse --short HEAD`" main.go

run:
	go run main.go

clean:
	rm image_server