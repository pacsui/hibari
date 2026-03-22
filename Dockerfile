FROM alpine:latest 

ARG TARGETARCH

COPY ./bin/hibari-${TARGETARCH} /usr/local/bin/hibari

RUN chmod +x /usr/local/bin/hibari

ENTRYPOINT ["/usr/local/bin/hibari"]