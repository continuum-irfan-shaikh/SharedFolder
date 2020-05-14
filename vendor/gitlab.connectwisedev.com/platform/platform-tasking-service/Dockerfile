# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8

# Copy the local package files to the container's workspace.
ADD . /go/src/gitlab.connectwisedev.com/platform/platform-tasking-service

# Changing working directory.
WORKDIR /go/src/gitlab.connectwisedev.com/platform/platform-tasking-service

# Building application.
RUN go build -o platform-tasking-service src/main.go

# Setting up environment variables

# Document that the service listens on port 12121.
EXPOSE 12121

# Run the platform-tasking-service command by default when the container starts.
CMD ["/go/src/gitlab.connectwisedev.com/platform/platform-tasking-service/platform-tasking-service"]
