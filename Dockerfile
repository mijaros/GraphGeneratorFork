FROM node:18-alpine as UI_BUILDER

COPY ui /ui/

WORKDIR /ui
RUN npm install && npm run build

FROM golang:1.20-alpine as GO_BUILDER

COPY ./ /src/
WORKDIR /src

RUN go build ./cmd/...

FROM alpine

ENV GENERATOR_UI_DIR "/var/ui/webassets"
ENV GENERATOR_PORT 8080
ENV GENERATOR_BIND_ADDR 0.0.0.0
ENV GENERATOR_HOST example.com
ENV GENERATOR_DB_ROOT "/var/badger/db"


COPY --from=UI_BUILDER /ui/dist/graphGenerator/ $GENERATOR_UI_DIR

COPY --from=GO_BUILDER /src/generator /bin/generator

ENTRYPOINT /bin/generator
CMD []
