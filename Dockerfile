# Build image with Go installation
FROM golang:alpine as builder

# Install NodeJs
RUN apk update && \
    apk add nodejs-npm

# Set working directory
WORKDIR /app

# Copy project into Docker image
COPY . .

# Create build directories
RUN mkdir /build && \
    mkdir /build/static

# Build Go executable
RUN go build -o /build/main web/main.go

# Build React assets
RUN cd react-client && \
    npm install && \
    npm run build && \
    cp -R build/* /build/static

# Deployment image with minimal linux installation
FROM alpine:latest

# Expose port 8080
EXPOSE 8080

# Copy build directory from previous stage
COPY --from=builder /build /app

# Set working directory
WORKDIR /app

# Create new user
RUN adduser -S -D -H -h /app appuser

# Change user
USER appuser

# Set entrypoint
ENTRYPOINT [ "./main" ]
