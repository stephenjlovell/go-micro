FROM alpine:latest

RUN mkdir /app

COPY frontApp /app
COPY ./cmd/web/templates /cmd/web/templates

CMD ["/app/frontApp"]