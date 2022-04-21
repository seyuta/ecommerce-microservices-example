FROM golang:1.15-alpine AS builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ .
RUN go build -o .build/s-catalog ./s-catalog/cmd/main.go

FROM alpine:latest
ARG server=default
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
ENV WORKDIR=/s-catalog
WORKDIR ${WORKDIR}
RUN mkdir ${WORKDIR}/log
COPY --from=builder /go/src/app/.build/s-catalog ${WORKDIR}/
COPY --from=builder /go/src/app/s-catalog/configs/${server}-config.yml ${WORKDIR}/configs/

CMD ./s-catalog ${server}