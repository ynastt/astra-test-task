FROM golang:1.20-alpine

# Work directory
WORKDIR /app

# Installing dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

# Copying all the files
COPY . ./

# RUN go build -o bin main.go
RUN GOOS=linux GOARCH=amd64 go build -o /astra-test main.go

CMD ["/astra-test", "report.sarif"]
