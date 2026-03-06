FROM debian:stable

# Install development toolchain
RUN apt update && apt install -y  \
    ca-certificates               \
    vim                           \
    git                           \
    curl                          \
    wget                          \
    procps                        \
    golang-1.24

# Setup environment variables
RUN echo ''                                         >> ~/.bashrc &&\
    echo 'export PATH="$PATH:/usr/lib/go-1.24/bin"' >> ~/.bashrc