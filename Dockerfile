FROM alpine

WORKDIR /app

COPY bin/hibari /app/hibari

CMD ["/app/hibari"]
