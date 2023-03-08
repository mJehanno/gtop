FROM tsl0922/ttyd:alpine as ttyd
FROM debian:stable-slim

RUN apt update && apt install -y curl chromium gpg wget git libnss3 ffmpeg bash libatk1.0-0 libatk-bridge2.0-0 libcups2 libxcomposite1 libxdamage1 sudo && useradd -G sudo --create-home vhs && passwd -d vhs

# Install VHS
RUN mkdir -p /etc/apt/keyrings && curl -fsSL https://repo.charm.sh/apt/gpg.key | gpg --dearmor -o /etc/apt/keyrings/charm.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | tee /etc/apt/sources.list.d/charm.list
RUN apt update && apt install -y vhs  
#&& echo 'kernel.unprivileged_userns_clone=1' > /etc/sysctl.d/userns.conf
# Install ttyd (vhs deps)
COPY --from=ttyd /usr/bin/ttyd /usr/bin/ttyd

USER vhs
WORKDIR /home/vhs
# Install Golang
RUN wget -q https://go.dev/dl/go1.20.2.linux-amd64.tar.gz && sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz && mkdir -p /home/vhs/go /home/vhs/app
ENV PATH=$PATH:/usr/local/go/bin
ARG GOBIN=/home/vhs/go
ARG PATH=$GOBIN:$PATH

# Install app
COPY --chown=vhs . /home/vhs/app
WORKDIR /home/vhs/app
RUN go install .
