FROM golang:latest AS build
#===================================
# Build the application from source
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o http-server

#===============================================
# Deploy the application binary into a lean image
FROM alpine:latest
COPY --from=build /go/src/app/http-server /goapp/http-server
WORKDIR /goapp
COPY . /throwaway
RUN cp -r /throwaway/templates ./templates || echo "No templates to copy"
RUN rm -rf /throwaway
RUN apk --no-cache add ca-certificates

EXPOSE 3000
#===================================
CMD  ["/goapp/http-server"]