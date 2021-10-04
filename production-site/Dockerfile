# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

# make working dir, use relative paths based on this directory
WORKDIR /app

# install basic compiler
RUN apk add --no-cache build-base
# copy and install all modules necessary to compile
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# copy other files
COPY *.go ./

# create the binary
RUN go build -o /production-site

# run 
#CMD [ "/production-site" ]

# test command
CMD [ "go", "test", "-run", "Test" ]