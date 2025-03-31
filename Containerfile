FROM golang:1.24 as builder
WORKDIR /build
COPY . .

ENV CGO_ENABLED=0
RUN make

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /build/bin/ .