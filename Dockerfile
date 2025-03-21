FROM golang:1.24-bookworm AS builder
WORKDIR /app
RUN useradd -r nonroot
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./bin ./cmd/api/

FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder --chown=nonroot:nonroot /app/bin/api .
USER nonroot
EXPOSE 8080
CMD ["/api"]


