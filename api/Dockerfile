FROM golang:1.22-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Add the directories which contain the golang scripts
COPY . .
COPY config.prod.yml ./config.yml

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -C cmd -o /api

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/api"]