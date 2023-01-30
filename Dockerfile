FROM golang:1.19.3-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN go build -o eks-test-app main.go


FROM gcr.io/distroless/static AS app
COPY --from=builder /app/eks-test-app /bin/eks-test-app

# Run as UID for nobody
USER 65534

ENTRYPOINT ["/bin/eks-test-app"]
