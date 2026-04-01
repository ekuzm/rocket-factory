FROM golang:1.26-alpine AS build

WORKDIR /app

COPY ./go.work ./go.work.sum ./
COPY ./order/go.mod ./order/go.sum ./order/
COPY ./inventory/go.mod ./inventory/go.sum ./inventory/
COPY ./payment/go.mod ./payment/go.sum ./payment/
COPY ./shared/go.mod ./shared/go.sum ./shared/
COPY ./platform/go.mod ./platform/go.sum ./platform/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./inventory/cmd/inventory-service/

FROM alpine:3.23

WORKDIR /app

COPY --from=build /app/inventory-service .

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    chown appuser:appgroup ./inventory-service

USER appuser

EXPOSE 50051

ENTRYPOINT [ "./inventory-service" ]
