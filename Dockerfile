FROM golang:1.22.7-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o subs-aggregator ./cmd/subs-aggregator/main.go

RUN CGO_ENABLED=0 go build -o migrate ./cmd/migrate/main.go

FROM alpine:latest

COPY --from=builder /app/subs-aggregator .
COPY --from=builder /app/migrate ./
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/${CONFIG_PATH} ./

EXPOSE 8080

CMD ./migrate up && ./subs-aggregator -c config.yaml
