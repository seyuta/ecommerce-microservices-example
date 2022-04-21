FROM golang:1.15-alpine AS builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ .
RUN go build -o .build/s-auth ./s-auth/cmd/main.go

FROM alpine:latest
ARG server=default
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
ENV WORKDIR=/s-auth
WORKDIR ${WORKDIR}
RUN mkdir ${WORKDIR}/log
COPY --from=builder /go/src/app/.build/s-auth ${WORKDIR}/
COPY --from=builder /go/src/app/s-auth/configs/${server}-config.yml ${WORKDIR}/configs/

CMD ./s-auth ${server}