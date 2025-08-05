FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates git

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cep-clima main.go

FROM alpine:latest AS final

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/cep-clima .

EXPOSE 8080

ENV WEATHER_API_KEY=4580cae74883437d87020027250408

COPY templates ./templates

CMD ["./cep-clima"]
