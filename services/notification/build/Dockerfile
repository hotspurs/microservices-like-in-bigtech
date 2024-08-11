FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/notification ./cmd/notification

FROM scratch AS final

WORKDIR /

COPY --from=build /bin/notification /notification

EXPOSE 8080

ENTRYPOINT ["/notification"]