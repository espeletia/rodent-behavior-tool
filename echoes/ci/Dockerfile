FROM alpine:3.20
ARG app="echoes"

COPY --from=mwader/static-ffmpeg:5.0.1-3 /ffmpeg /usr/bin/
COPY --from=mwader/static-ffmpeg:5.0.1-3 /ffprobe /usr/bin/

RUN apk update
RUN apk add graphicsmagick

ADD ./${app}/configurations /app/configurations
ADD ./${app}/videos /app/videos
ADD ./${app}/build /app
# ADD ./${app}/migrations /app/migrations

WORKDIR /app

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot
RUN chown -R nonroot:nonroot /app
USER nonroot

ENTRYPOINT ./start_server

