FROM golang:1.18 as build

RUN mkdir build

WORKDIR /workspace

COPY go.mod ./
COPY go.sum ./

RUN go mod download -json

COPY cmd ./cmd
COPY pkg ./pkg

RUN CGO_ENABLED=0 go build -o /api ./cmd/dbus-api
RUN go test ./...

FROM scratch

COPY --from=build /api /api
ENTRYPOINT ["/api"]
