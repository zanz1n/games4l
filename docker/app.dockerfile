FROM golang:1.20.4-alpine3.16 AS builder

WORKDIR /build

RUN apk add make

COPY . .

RUN make build-min

FROM alpine:3.16

COPY --from=builder /build/dist/stable.bin /app

CMD [ "/app" ]
