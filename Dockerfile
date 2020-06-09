FROM golang:1.14-alpine3.11 AS build

WORKDIR /app
COPY . .

RUN go build -o ./robin ./cmd/robin/main.go

FROM alpine:3.11

WORKDIR /app
COPY --from=build /app/robin ./robin

ENTRYPOINT [ "./robin" ]
