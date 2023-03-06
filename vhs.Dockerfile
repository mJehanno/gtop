FROM alpine:latest

RUN apk add go firefox ffmpeg ttyd bash chromium sudo git
RUN    	echo '%wheel ALL=(ALL) ALL' > /etc/sudoers.d/wheel; \
        adduser -D vhs -G wheel; 


USER vhs
RUN mkdir -p /home/vhs/go

ARG GOBIN=/home/vhs/go
ARG PATH=$GOBIN:$PATH
ADD . /app/
WORKDIR /app
RUN go install .

