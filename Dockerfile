FROM golang:1.17 AS server

RUN apt-get update && apt-get dist-upgrade

COPY server /build/server
WORKDIR /build/server
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata -o server

########################################

FROM node:16 AS client

RUN apt-get update && apt-get dist-upgrade

COPY client /build/client
WORKDIR /build/client
RUN yarn install
RUN yarn build

########################################

FROM busybox
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=server /build/server/server /app/server
COPY --from=client /build/client/dist/ /app/static/

EXPOSE 9081
CMD [ "./server" ]
