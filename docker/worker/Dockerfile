FROM golang:1.20
WORKDIR /app

COPY . .

RUN go install github.com/cosmtrek/air@v1.43.0
RUN go mod tidy

RUN apt-get update
RUN apt-get install -y ffmpeg imagemagick


CMD ["air", "-c", "./cmd/worker/.air.toml"]