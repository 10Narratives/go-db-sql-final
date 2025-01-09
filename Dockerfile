FROM golang:1.23

WORKDIR /parcel-store

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /parcel-app

CMD [ "/parcel-app" ]