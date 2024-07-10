FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cosmtrek/air@v1.51.0
RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN go build -o ./api/main ./api/main.go

EXPOSE 8081

CMD "air"
