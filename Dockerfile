FROM golang:1.26-alpine3.23 AS builder

# Set working directory inside builder
WORKDIR /app

# Cache go modules
COPY ./demo-api/go.mod ./demo-api/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary from the demo-api package by changing into that directory
WORKDIR /app/demo-api
RUN CGO_ENABLED=0 go build -o /app/out/api-server .

FROM scratch
ENV GIN_MODE=release
# Copy the compiled binary from the builder stage
COPY --from=builder /app/out/api-server /

# scratch has no certificates, but binary uses none
EXPOSE 8080

CMD ["./api-server"]
