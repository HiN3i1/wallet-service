FROM golang:1.16.7-alpine as builder

RUN apk add make git

ADD . /wallet-service
WORKDIR /wallet-service

RUN git config --global credential.helper "store --file `pwd`/.git-credentials"
RUN make all

FROM alpine:latest

COPY --from=builder /wallet-service/build/bin/* /usr/local/bin/
COPY .env.example .env

RUN apk --no-cache add ca-certificates

ENTRYPOINT [ "/usr/local/bin/wallet-service", "-m", "run-server" ]

EXPOSE 10445