FROM golang:1.20
WORKDIR /app

COPY . .

RUN go install github.com/cosmtrek/air@v1.43.0
RUN go mod tidy


CMD ["air", "-c", "./cmd/backend/.air.toml"]