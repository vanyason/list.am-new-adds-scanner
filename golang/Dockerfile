# Use Ubuntu as the base image
FROM ubuntu:20.04

# Set the working directory
WORKDIR /app

# Copy the source code to the working directory
COPY . .

# Install the required packages
RUN apt-get update && apt-get install -y git && \
apt-get install -y wget && \
apt-get install -y gcc

# Download and install the latest version of Golang
RUN wget https://dl.google.com/go/go1.19.linux-amd64.tar.gz && \
    tar -xvf go1.19.linux-amd64.tar.gz && \
    mv go /usr/local

# Set the environment variable for Go
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Go to the app directory
WORKDIR /app

# Download and install the dependencies defined in the go.mod file
RUN go mod download
