FROM golang:1.22-alpine AS build
RUN apk --no-cache --update add build-base
ADD . /app/
WORKDIR /app

RUN go build -o api .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["/app/api"]