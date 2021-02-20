
# --- buider stage

FROM golang:1.15-buster AS builder

# Set working directory
WORKDIR /src

# Prepare Golang environment
ENV GO111MODULE on

# Copy go-modules definition for the application
COPY go.mod go.mod

# Download required go-modules
RUN go mod download

# Copy application code
COPY . .

# Compile application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/app


# --- final stage

FROM alpine

# Install additionals
RUN apk add --no-cache bash iputils curl bind-tools

# Copy application executable
WORKDIR /bin/
COPY --from=builder --chown=1001 /bin/app .

## Metrics
EXPOSE 9090
## KubeProbes
# EXPOSE 9091

# Set user
USER 1001

# Run application
ENTRYPOINT "/bin/app"
