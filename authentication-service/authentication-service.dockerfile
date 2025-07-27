# Dockerfile for the authentication service
# This Dockerfile is based on Alpine Linux

FROM alpine:latest


# Create a directory for the application
RUN mkdir /app


# Copy the authentication app binary to the /app directory
COPY authApp /app


# Set the default command to run the authentication app
CMD [ "/app/authApp"]