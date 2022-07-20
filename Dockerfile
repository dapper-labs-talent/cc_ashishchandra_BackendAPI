FROM golang:latest
WORKDIR /app
COPY . .
RUN go mod download
RUN go build .
CMD ./cc_ashishchandra_BackendAPI
