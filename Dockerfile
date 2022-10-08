FROM alpine

# Copy the built binary over into the image
COPY build/scraparr /scraparr

# Install common CA certificates, otherwise HTTPS requests return unknown authority errors
RUN apk add --no-cache ca-certificates

# Expose the port the server listens on
EXPOSE 8080/tcp

# Start the server
ENTRYPOINT [ "/scraparr" ]