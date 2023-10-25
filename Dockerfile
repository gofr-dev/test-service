FROM alpine:edge

WORKDIR /src

COPY ./configs /configs

COPY build/test-service ./test-service
RUN chmod +x ./test-service

EXPOSE 80

ENTRYPOINT [ "./test-service" ]