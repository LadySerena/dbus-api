FROM golang:1.14

RUN mkdir build

WORKDIR /workspace

COPY go.mod ./
COPY go.sum ./

RUN go mod download -json

COPY cmd ./cmd
COPY pkg ./pkg

RUN go build -o ./build ./cmd/dbus-api
RUN go test ./...