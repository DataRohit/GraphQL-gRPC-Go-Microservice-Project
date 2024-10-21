FROM golang:1.23.2-bullseye AS build-stage

WORKDIR /app

COPY go.work go.work.sum /app/
COPY account/go.mod account/go.sum /app/account/
COPY . .

WORKDIR /app/account
RUN go mod download

WORKDIR /app/account
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/account/bin/main /app/account/cmd/

FROM golang:1.23.2-alpine AS release-stage

RUN go install github.com/air-verse/air@latest
RUN apk update && apk add --no-cache curl

WORKDIR /app/account

COPY --from=build-stage /app/account/bin/main /app/account/bin/main
COPY . .

CMD ["air"]
