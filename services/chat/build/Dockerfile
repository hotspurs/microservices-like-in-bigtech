FROM golang:1.22.5-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/chat ./cmd/chat

FROM scratch AS final

WORKDIR /

COPY --from=build /bin/chat /chat
COPY swagger swagger

EXPOSE 8080
EXPOSE 8082

ENTRYPOINT ["/chat"]