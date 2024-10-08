FROM golang:1.21-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . ./

# Set necessary environment variables needed
# for our image and build the sender.
RUN go build -o storage cmd/main.go

EXPOSE 8080
ENTRYPOINT ["./storage"]