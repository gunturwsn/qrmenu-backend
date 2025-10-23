# build stage
FROM golang:1.23.5-alpine AS builder
ENV CGO_ENABLED=0
ENV GOTOOLCHAIN=auto
WORKDIR /app
RUN apk add --no-cache build-base git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o api ./cmd/api

# runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/api /usr/local/bin/api
COPY migrations /app/migrations
EXPOSE 8080
ENTRYPOINT ["api"]
