FROM golang:1.23 AS build
WORKDIR /app
RUN useradd -u 1001 nonroot
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o /rss-aggregator

FROM alpine
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /app/rss-aggregator rss-aggregator
USER nonroot
EXPOSE 8080
CMD ["./rss-aggregator"]



