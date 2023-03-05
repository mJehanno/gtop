FROM alpine:latest

RUN apk add go firefox ffmpeg ttyd bash chromium sudo git
RUN    	echo '%wheel ALL=(ALL) ALL' > /etc/sudoers.d/wheel; \
        adduser -D vhs -G wheel; 

RUN mkdir -p /home/docker/go
ARG GOBIN=/home/docker/go
ARG PATH=$GOBIN:$PATH

ADD . /app/
WORKDIR /app
RUN go install .
RUN go install github.com/charmbracelet/vhs@latest

USER vhs
