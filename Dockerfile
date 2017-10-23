# Start by building the application.
FROM golang:1.8 as build

WORKDIR /go/src/app
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
# RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /go/bin/app /
CMD ["/app"]