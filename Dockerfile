#####################
## Go plugins
#####################
FROM kong:2.2.0-ubuntu as compiler

# Install build tools
USER root
RUN apt-get update -y && apt-get install --no-install-recommends -y -q curl build-essential ca-certificates git

# Download and configure Go compiler
RUN curl -s https://dl.google.com/go/go1.15.3.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH $PATH:/usr/local/go/bin

# Copy and compile go-plugins
RUN go get github.com/Masterminds/sprig
RUN go get github.com/Kong/go-pluginserver

#RUN git clone https://github.com/Kong/go-plugins /usr/src/go-plugins
COPY . /usr/src/go-plugins
RUN mkdir /tmp/go-plugins; cp /usr/src/go-plugins/midnight.go /tmp/go-plugins
RUN cd /tmp/go-plugins; go build -buildmode plugin midnight.go

#####################
## Release image
#####################
FROM kong:2.2.0-ubuntu

COPY --from=compiler /tmp/go-plugins/*.so /usr/local/kong/
COPY --from=compiler /go/bin/go-pluginserver /usr/local/bin/go-pluginserver

USER kong