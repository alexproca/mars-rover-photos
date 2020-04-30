FROM alpine:3.11

RUN mkdir /app

WORKDIR /app

COPY server /app
COPY .env /app
COPY templates /app/templates
COPY public /app/public

ENV INTERFACE=0.0.0.0

CMD ["/app/server"]
