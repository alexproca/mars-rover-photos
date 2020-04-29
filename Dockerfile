FROM alpine:3.11

RUN mkdir /app

WORKDIR /app

COPY server /app
COPY .env /app
COPY templates /app/templates
COPY public /app/public

CMD ["/app/server"]
