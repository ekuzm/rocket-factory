FROM golang:1.26 AS builder

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

FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app

COPY --from=builder /app/inventory-service .

EXPOSE 50051

ENTRYPOINT [ "./inventory-service" ]
