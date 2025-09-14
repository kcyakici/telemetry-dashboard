FROM golang:1.25.1-alpine3.22 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-output

# FROM build-stage AS run-test-stage
# RUN go test -v ./...

FROM alpine:3.22.1 AS build-release-stage
WORKDIR /
COPY --from=build-stage /go-output /go-output
EXPOSE 8080
ENTRYPOINT [ "/go-output" ]