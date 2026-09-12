FROM golang:1.25-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/scheduler \
    .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /out/scheduler /app/scheduler
COPY web /app/web

RUN mkdir -p /data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

EXPOSE 7540

CMD ["/app/scheduler"]