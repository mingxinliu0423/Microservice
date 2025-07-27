# base go image
FROM golang:1.23-alpine as builder


# Create a directory for the app
RUN mkdir /app


# Copy the app source code into the container
COPY . /app


# Set the working directory to the app directory
WORKDIR /app


# Build the app
RUN CGO_ENABLED=0 go build -o brokerApp ./service


# Make the binary executable
RUN chmod +x /app/brokerApp


# Build a tiny docker image
FROM alpine:latest


# Create a directory for the app
RUN mkdir /app


# Copy the built binary from the builder stage to the final image
COPY --from=builder /app/brokerApp /app


# Specify the command to run when the container starts
CMD [ "/app/brokerApp" ]
