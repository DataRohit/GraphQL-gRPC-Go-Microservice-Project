FROM golang:1.23.2-bullseye AS build-stage

WORKDIR /app

COPY go.work go.work.sum /app/
COPY account/go.mod account/go.sum /app/account/
COPY gateway/go.mod gateway/go.sum /app/gateway/
COPY . .

WORKDIR /app/account
RUN go mod download

WORKDIR /app/gateway
RUN go mod download

WORKDIR /app/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gateway/bin/main /app/gateway/

FROM golang:1.23.2-alpine AS release-stage

RUN go install github.com/air-verse/air@latest
RUN apk update && apk add --no-cache curl

WORKDIR /app/gateway

COPY --from=build-stage /app/gateway/bin/main /app/gateway/bin/main
COPY . .

CMD ["air"]
