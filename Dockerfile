FROM golang

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .

RUN go build -o . /app/cmd/main.go
EXPOSE 8080
RUN chmod +x ./main
CMD ["./main"]