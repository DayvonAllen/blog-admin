FROM golang:1.16.3-alpine

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/blog-admin-server

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/blog-admin-server .

# This container exposes port 8082 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./out/blog-admin-server"]
