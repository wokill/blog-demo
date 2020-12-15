.PHONY: build gdoc doc run publish clean lint all

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

repo=registry-in.dustess.com:9000

project_name=mk-blog-svc

imageName=markting/${project_name}
tag=`git branch | grep \* | cut -d ' ' -f2`
imageWholeName=${repo}/${imageName}:${tag}
port=5000
grpc_port=50000

build:
	go mod vendor
	go build -o ${project_name} cmd/server/main.go
	docker build -t ${imageWholeName} .

gdoc:
	rm -f api/swagger/* || true
	swag init -g internal/http/server.go -o api/swagger

doc:
	mkdir docs || true
	rm -f docs/* || true
	swag init -g internal/http/server.go -o api/swagger
	sleep 1
	yapi import || true

publish:
	docker push ${imageWholeName}

run:
	docker run --name ${imageName} -p ${port}:${port} -p ${grpc_port}:${grpc_port} -v server.conf:/app/server.conf ${imageWholeName}

clean:
	docker rm -f ${imageName}

lint:
	go list ./... | grep -v /vendor/ | xargs -L1 golint

update:
	go get git.dustess.com/mk-base/log@master
	go get git.dustess.com/mk-base/es-driver@master
	go get git.dustess.com/mk-base/gin-ext@master
	go get git.dustess.com/mk-base/mongo-driver@master
	go get git.dustess.com/mk-base/redis-driver@master
	go get git.dustess.com/mk-base/oss-driver@master
	go mod tidy

all: build run
