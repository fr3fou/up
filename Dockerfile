FROM golang:latest
COPY . .
RUN go build -o up ./
EXPOSE 8080
CMD ['./main']
