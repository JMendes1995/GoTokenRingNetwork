FROM debian:trixie
WORKDIR /usr/src/app

RUN apt-get update -y
RUN apt-get install -y \
  build-essential \
  curl \
  net-tools \
  pkg-config \
  libssl-dev \
  vim \
  iputils-ping \
  ssh \
  golang


COPY . .
RUN go mod download && go mod verify
