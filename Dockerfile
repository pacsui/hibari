FROM alpine:latest 

ARG TARGETARCH

COPY ./bin/hibari-${TARGETARCH} /bin/hibari

RUN chmod +x /bin/hibari

WORKDIR /app

ENTRYPOINT ["hibari"]
