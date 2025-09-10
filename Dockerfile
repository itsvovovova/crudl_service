FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install gotest.tools/gotestsum@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init -g src/cmd/main.go -o ./docs
RUN go build -o project ./src/cmd

EXPOSE 8080

CMD ["./project"]