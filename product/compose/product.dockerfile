FROM golang:1.23.2-bullseye AS build-stage

WORKDIR /app

COPY go.work go.work.sum /app/
COPY product/go.mod product/go.sum /app/product/
COPY . .

WORKDIR /app/product
RUN go mod download

WORKDIR /app/product
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/product/bin/main /app/product/cmd/

FROM golang:1.23.2-alpine AS release-stage

RUN go install github.com/air-verse/air@latest
RUN apk update && apk add --no-cache curl

WORKDIR /app/product

COPY --from=build-stage /app/product/bin/main /app/product/bin/main
COPY . .

CMD ["air"]
