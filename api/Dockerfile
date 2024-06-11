

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

#Stage 2: Run tests
FROM golang:1.22-alpine as test
WORKDIR /app
COPY . .
CMD ["go", "test", "./tests"]

#Stage 3: Prepare release
FROM alpine:3.20 as prepare
#Download and microsoft CA Certiciates
RUN apk update && apk add curl
WORKDIR /usr/local/share/ca-certificates
RUN curl -ks 'https://cacerts.digicert.com/BaltimoreCyberTrustRoot.crt' -o '/usr/local/share/ca-certificates/BaltimoreCyberTrustRoot.crt'
RUN curl -ks 'https://cacerts.digicert.com/DigiCertGlobalRootCA.crt' -o '/usr/local/share/ca-certificates/DigiCertGlobalRootCA.crt'
RUN curl -ks 'https://cacerts.digicert.com/DigiCertGlobalRootG2.crt' -o '/usr/local/share/ca-certificates/DigiCertGlobalRootG2.crt'
RUN curl -ks 'https://cacerts.digicert.com/DigiCertGlobalRootG3.crt' -o '/usr/local/share/ca-certificates/DigiCertGlobalRootG2.crt'
RUN curl -ks 'https://www.microsoft.com/pkiops/certs/Microsoft%20ECC%20Root%20Certificate%20Authority%202017.crt' -o '/usr/local/share/ca-certificates/MicrosoftRootCertificateAuthority.crt'
#Add microsoft CA to trusted root certificates
RUN /usr/sbin/update-ca-certificates

#Stage $: Run Api
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