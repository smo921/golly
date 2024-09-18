FROM golang:latest AS builder
WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./golly


FROM scratch
WORKDIR /app
COPY --from=builder /build/golly ./golly
CMD ["./golly"]
