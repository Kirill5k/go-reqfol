FROM golang:1.21-alpine as dependencies

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM dependencies as build
COPY . ./
RUN CGO_ENABLED=0 go build -o /main -ldflags="-w -s" .

FROM golang:1.21-alpine

ARG VERSION
ENV VERSION=$VERSION

COPY --from=build /main /main
COPY --from=build /app/config/application.yaml /go/config/application.yaml
CMD ["/main"]
