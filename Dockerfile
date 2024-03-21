FROM gocv/opencv:4.8.1
COPY ./ /app
WORKDIR /app
RUN go build -ldflags="-w -s" -o /binary

ENTRYPOINT ["/binary"]
