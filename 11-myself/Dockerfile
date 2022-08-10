FROM golang:1.19 as builder

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN CGO_ENABLE=0 go build -o ingress-manager main.go

FROM alpine:3.15.3

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk --no-cache add tzdata ca-certificates libc6-compat libgcc libstdc++ alsa-lib-dev

WORKDIR /app

COPY --from=builder /app/ingress-manager .

CMD ["./ingress-manager"]

