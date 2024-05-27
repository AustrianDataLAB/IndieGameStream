#Stage 1: Compile and build
FROM golang:1.22-alpine as build
# Set destination for COPY
WORKDIR /app
# Download Go modules
COPY go.mod go.sum ./
RUN go mod download
# Add the directories which contain the golang scripts
COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -C cmd -o /api

#Stage 2a: Run tests
FROM golang:1.22-alpine as test
WORKDIR /app
COPY . .
CMD ["go", "test", "./tests"]

#Stage 2b: Run Api
FROM scratch as release
#Copy migration scripts
COPY migrations migrations
#Copy build result to next stage
COPY --from=build /api /api
EXPOSE 8080
# Run
CMD ["/api"]