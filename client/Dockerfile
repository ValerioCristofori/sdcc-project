# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

# make working dir, use relative paths based on this directory
WORKDIR /app

# copy and install all modules necessary to compile
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy other files
COPY *.go ./

# create the binary
RUN go build -o /client

# run 
CMD [ "/client" ]