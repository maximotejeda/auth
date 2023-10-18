ARG GOIMAGE=bookworm

FROM golang:${GOIMAGE} AS build
WORKDIR /app
COPY go.mod go.sum .
# those mount are saved to be reused on all steps needed
RUN --mount=type=cache,target=/go/pkg/mod/ \
    go mod download -x
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod/ \
    CGO_ENABLED=0 go build -o /bin/auth cmd/auth/auth.go

FROM scratch AS server
COPY --from=build /bin/auth /bin/

ENTRYPOINT ["/bin/auth"]
