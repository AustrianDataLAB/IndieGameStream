#Stage 1: Compile and build
FROM golang:1.20.14-alpine3.18 as build
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

#Stage 3: Prepare release
FROM alpine:3.20 as prepare
#Download Google CA Certiciates
RUN apk update && apk add curl
WORKDIR /usr/local/share/ca-certificates
RUN curl -ks 'https://i.pki.goog/r1.pem' -o '/usr/local/share/ca-certificates/r1.pem'
RUN curl -ks 'https://i.pki.goog/r2.pem' -o '/usr/local/share/ca-certificates/r2.pem'
RUN curl -ks 'https://i.pki.goog/r3.pem' -o '/usr/local/share/ca-certificates/r3.pem'
RUN curl -ks 'https://i.pki.goog/r4.pem' -o '/usr/local/share/ca-certificates/r4.pem'
RUN curl -ks 'https://i.pki.goog/gsr4.pem' -o '/usr/local/share/ca-certificates/gsr4.pem'
#Add Google CA to trusted root certificates
RUN /usr/sbin/update-ca-certificates


#Stage 2b: Run Api
FROM scratch as release
#Copy trusted CA certificates
COPY --from=prepare /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
#Copy migration scripts
COPY migrations migrations
#Copy build result to next stage
COPY --from=build /api /api
EXPOSE 8080
# Run
CMD ["/api"]