# Specifies a parent image
FROM golang:1.20.2-bullseye

LABEL org.opencontainers.image.source=https://github.com/EdisonAltamirano/security_backend




RUN apt-get update && apt-get install -y  \
    libgstreamer1.0-dev  \
    libgstreamer-plugins-base1.0-dev  \
    libgstreamer-plugins-bad1.0-dev  \
    gstreamer1.0-plugins-base  \
    gstreamer1.0-plugins-good  \
    gstreamer1.0-plugins-bad  \
    gstreamer1.0-plugins-ugly  \
    gstreamer1.0-libav  \
    gstreamer1.0-tools  \
    gstreamer1.0-x  \
    gstreamer1.0-alsa  \
    gstreamer1.0-gl  \
    gstreamer1.0-gtk3  \
    gstreamer1.0-qt5  \
    gstreamer1.0-pulseaudio

# Creates an app directory to hold your app’s source code
WORKDIR /camera_streamer

COPY go.mod .
COPY go.sum .

# Installs Go dependencies
RUN go mod download

# Copies everything from your root directory into /app
COPY . .

# Builds your app with optional configuration
RUN go build -buildvcs=false -o ./camera_streamer github.com/EdisonAltamirano/security_backend.git

ENV CAMERA_SERVER_CONFIG=/config

# Tells Docker which network port your container listens on
EXPOSE 3000

# Specifies the executable command that runs when the container starts
ENTRYPOINT [ "/camera_streamer/camera_streamer" ]