FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8081

RUN go build -o ./cmd/api/main ./cmd/api/main.go
RUN go build -o ./cmd/worker/production_failed_worker/main ./cmd/worker/production_failed_worker/main.go

CMD [ "./cmd/api/main" ]
