FROM golang:1.23 AS server-builder

WORKDIR /build

ENV CGO_ENABLED=0

COPY . .

RUN go build -o /socks5-server

# final image

FROM ubuntu

COPY --from=server-builder /socks5-server /socks5-server
