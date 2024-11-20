FROM alpine:3.20
ARG app="tusk"
ARG s3Url="http://minio:9000"
ENV S3_URL=$s3Url

ADD ./${app}/configurations /app/configurations
ADD ./${app}/build /app
ADD ./${app}/migrations /app/migrations

WORKDIR /app


RUN addgroup -S nonroot -g 2000 \
    && adduser -S nonroot -G nonroot -u 2000
RUN chown -R nonroot:nonroot /app
USER nonroot

ENTRYPOINT ./start_server

