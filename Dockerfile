FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/orbit ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/orbit /app/orbit
COPY docker/entrypoint.sh /entrypoint.sh

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/orbit"]
