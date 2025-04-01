FROM alpine:latest

RUN mkdir /app

COPY . /app
COPY internal/database/migrations /app/internal/database/migrations

WORKDIR /app

CMD ["/app/providerApp"]
