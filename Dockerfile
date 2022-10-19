FROM golang:1.19 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o detek ./main.go

FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/detek /detek
USER nonroot:nonroot
ENTRYPOINT ["/detek"]
CMD [ "run" ]
