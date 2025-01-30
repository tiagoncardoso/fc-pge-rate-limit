FROM golang:1.23 AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rl-app ./cmd/app-rl/main.go

FROM golang:1.23 AS prod

WORKDIR /app
COPY --from=build /app/rl-app .
COPY --from=build /app/.env .env

EXPOSE 8080

ENTRYPOINT ["./rl-app"]