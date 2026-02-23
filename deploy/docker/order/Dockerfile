FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY ./go.work ./go.work.sum ./
COPY ./order/go.mod ./order/go.sum ./order/
COPY ./inventory/go.mod ./inventory/go.sum ./inventory/
COPY ./payment/go.mod ./payment/go.sum ./payment/
COPY ./shared/go.mod ./shared/go.sum ./shared/

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build ./order/cmd/order-service/

FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/order-service .

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown appuser:appgroup ./order-service

USER appuser

EXPOSE 80

ENTRYPOINT [ "./order-service" ]
