#BUILD
FROM golang:1.13-alpine AS build
RUN mkdir /app/
WORKDIR /app/
COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s" -o raft 

FROM scratch AS run
COPY --from=build /app/raft .
ENTRYPOINT [ "/raft"]