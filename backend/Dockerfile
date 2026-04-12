FROM golang:1.25-alpine

WORKDIR /app

# install git (needed sometimes)
RUN apk add --no-cache git

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy project
COPY . .

# build binary
RUN go build -o main ./cmd

# run app
CMD ["./main"]