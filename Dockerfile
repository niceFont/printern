
#build stage
FROM golang:alpine AS builder


WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v github.com/julienschmidt/httprouter
#RUN go install -v ./...


EXPOSE 3000

CMD [ "go run", "./cmd/printern.go" ]