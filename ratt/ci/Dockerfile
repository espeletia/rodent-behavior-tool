FROM ultralytics/ultralytics:8.2.102

# ENV DEBIAN_FRONTEND=noninteractive
COPY --from=mwader/static-ffmpeg:5.0.1-3 /ffmpeg /usr/bin/
COPY --from=mwader/static-ffmpeg:5.0.1-3 /ffprobe /usr/bin/

RUN apt-get update
RUN apt-get install -y graphicsmagick


# RUN apt-get update && apt-get install -y \
#     ffmpeg \
#     libsm6 \
#     libxext6 \
#     libdvdnav4 \
#     gstreamer1.0-plugins-bad \
#     gstreamer1.0-plugins-ugly \
#     libdvd-pkg
#
# RUN dpkg-reconfigure libdvd-pkg

# COPY --from=mwader/static-ffmpeg:5.0.1-3 /ffmpeg /usr/bin/
# COPY --from=mwader/static-ffmpeg:5.0.1-3 /ffprobe /usr/bin/
#
# RUN apt-get update -y
# RUN apt-get install graphicsmagick -y


# RUN apt install ubuntu-restricted-extras -y

# RUN apt-get update && apt-get install -y \
#     libx264-dev \
#     ffmpeg \
#     && apt-get clean && rm -rf /var/lib/apt/lists/*
RUN ffmpeg -encoders
RUN ffmpeg -encoders
RUN ffmpeg -encoders
RUN ffmpeg -encoders


WORKDIR /app
ADD ./ratt/requirements.txt /app/requirements.txt
RUN pip3 install -r requirements.txt

ADD ./ratt/ /app/

ENTRYPOINT ["python", "main.py"]

