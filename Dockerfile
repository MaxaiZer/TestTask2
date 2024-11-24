FROM scratch AS base

FROM golang:1.23 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY docs/ ./docs/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/build/main cmd/main.go

FROM base AS final
WORKDIR /app
COPY configs ./configs/
COPY migrations ./migrations/
COPY --from=build /app/build/main .
CMD ["./main"]