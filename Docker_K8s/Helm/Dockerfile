FROM golang:1.19.6 AS build_base
RUN apt-get install -y git
# Set the Current Working Directory inside the container
RUN mkdir /tmp/diaginfra
WORKDIR /tmp/diaginfra
# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
# Unit tests
RUN go test ./... -v
# Build the Go app
RUN CGO_ENABLED=0 go build -o ./out/diaginfra -buildvcs=false main.go

FROM node:18.15.0-alpine AS nodejs
WORKDIR /tmp/app
COPY . .
RUN npm install
RUN npm run build

# Start fresh from a smaller image
FROM alpine:3.9 
COPY --from=build_base /tmp/diaginfra/out/diaginfra /app/
COPY --from=nodejs /tmp/app/build /app/build
COPY ./conf/config.yaml /app/conf/
# This container exposes port 3000 to the outside world
EXPOSE 3000
# Run the binary program produced by `go install`
WORKDIR /app
CMD ["diaginfra" "server"]