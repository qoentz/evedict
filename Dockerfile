# Use the official Go image for the build
FROM golang:1.21.1-alpine

# Set the working directory inside the container to /evedict
WORKDIR /evedict

# Copy the entire project into the container
COPY . .

# Download dependencies
RUN go mod download

# Build the Go app from the /evedict/cmd directory and place the output binary in /evedict
RUN go build -o /evedict/main /evedict/cmd/app/main.go

# Expose port 8080
EXPOSE 8080

# Run the compiled binary from /evedict
CMD ["/evedict/main"]




