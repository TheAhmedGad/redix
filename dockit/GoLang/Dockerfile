FROM golang:1.22-alpine

WORKDIR /app

# Install build dependencies for cgo/sqlite3
RUN apk add --no-cache gcc musl-dev sqlite-dev

# (Optional) Install air for live reload
RUN go install github.com/cosmtrek/air@v1.40.4

ENV PATH="/root/go/bin:$PATH"

# Enable CGO for go test
ENV CGO_ENABLED=1

CMD ["air"]