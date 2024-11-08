FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache make

COPY . .

RUN make build_server

FROM scratch

COPY --from=builder /app/bin/gophkeeper-server /server

ENTRYPOINT ["/server"]
