FROM alpine:latest

RUN mkdir /app

COPY mockApp /app

CMD ["/app/mockApp"]

