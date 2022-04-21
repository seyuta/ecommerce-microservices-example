FROM golang:1.15-alpine AS builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ .
RUN go build -o .build/transactions ./s-transactions/cmd/main.go

FROM alpine:latest
ARG server=default
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
ENV WORKDIR=/transactions
WORKDIR ${WORKDIR}
RUN mkdir ${WORKDIR}/log
COPY --from=builder /go/src/app/.build/transactions ${WORKDIR}/
COPY --from=builder /go/src/app/s-transactions/configs/${server}-config.yml ${WORKDIR}/configs/

CMD ./transactions ${server}