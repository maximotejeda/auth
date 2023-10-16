.PHONY: test clean build run-local run-build
NAME="auth"
OSS="linux windows darwin"
PLATFORM="amd64 arm64"
ACTUALOS=`uname | tr '[:upper:]' '[:lower:]'`
ACTUALPLATFORM=`uname -r | sed 's/\([[:digit:].]*-\)//g'`
ACTUALBIN="${NAME}-${ACTUALOS}-${ACTUALPLATFORM}"
TEMPFILES="key* *.db bin"

build-image:
build-image-debug:
run-image:
run-image-debug:
build-all:
	@echo "building bin on ./bin/${ALLGOFILES}"
	@for j in "${OSS}"; do \
	for i in "${PLATFORM}"; do \
	echo "building => bin/${NAME}-$$j-$$i" && \
	GOOS="$$j" GOARCH="$$i" go build -o bin/"${NAME}"-"$$j"-"$$i" cmd/auth/auth.go ; done \
	done
build:
	@echo  "building bin/${ACTUALBIN}"
	@GOOS=${ACTUALOS} GOARCH=${ACTUALPLATFORM} go build -o bin/${ACTUALBIN} cmd/auth/auth.go
run-build: build
	@echo  "running bin/${NAME}-${ACTUALOS}-${ACTUALPLATFORM}"
	@bin/"${ACTUALBIN}"
run-local:
	@echo "runnning code" 
	@go run cmd/auth/auth.go -timeout=9s -addr=:8080 -debug=1
test:
	@echo "running test"
	@go test -v ./...
clean:
	@echo "deleting keys binaries and database"
	@rm -rf key* *.db bin
