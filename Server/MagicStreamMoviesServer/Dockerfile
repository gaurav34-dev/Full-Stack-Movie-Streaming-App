FROM golang:1.20-alpine AS build
WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /server/main ./

FROM alpine:latest
RUN adduser -D -u 1000 appuser
COPY --from=build /server/main /server/main
USER appuser
EXPOSE 8081
CMD ["/server/main"]
