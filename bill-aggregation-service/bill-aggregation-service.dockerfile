FROM alpine:latest

RUN mkdir /app

COPY billApp /app

CMD ["/app/billApp"]
