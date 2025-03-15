# Build the application from source
FROM golang:1.24rc3-bookworm AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY config/ config/
COPY domain/ domain/
COPY infrastructure/ infrastructure/
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /web

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /web /web

# EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/web"]
