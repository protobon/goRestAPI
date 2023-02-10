FROM golang:alpine

# Add Maintainer info
LABEL maintainer="Ayrton Coelho"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

#Setup hot-reload for dev stage
# RUN go get github.com/githubnemo/CompileDaemon
# RUN go get -v golang.org/x/tools/gopls

# ENTRYPOINT CompileDaemon --build="go build -a -installsuffix cgo -o main ." --command=./main

# build app and expose port 8080
RUN go build -o /api
EXPOSE 8080

# run application
CMD [ "/api" ]