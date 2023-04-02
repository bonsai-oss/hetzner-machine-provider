FROM golang:1.20 as builder
WORKDIR /build
COPY . .

ENV CGO_ENABLED=0
RUN make

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /build/bin/ .