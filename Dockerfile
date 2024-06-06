FROM golang:latest as depedencies
LABEL authors="fariz.wisnu"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download


FROM depedencies as builder
COPY . ./
RUN CGO_ENABLED=0 go build -o /main -ldflags "-s -w" .

FROM golang:latest
COPY --from=builder /main /main
CMD ["/main"]