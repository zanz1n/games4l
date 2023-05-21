FROM golang:1.20.4-alpine3.18 AS builder

WORKDIR /build

RUN apk add make

COPY ./go.work ./go.work.sum /build/

COPY ./libs /build/libs

COPY ./services /build/services

WORKDIR /build/services/telemetry

RUN make build

FROM gcr.io/distroless/base-debian11

COPY --from=builder /build/services/telemetry/dist/stable.bin /app

CMD [ "/app" ]
