FROM kubemq/gobuilder as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq/demo/backend/elastic
WORKDIR $GOPATH/github.com/kubemq/demo/backend/elastic
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -o kubemq-elastic-run .
FROM alpine:latest
RUN apk add ca-certificates
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
COPY --from=builder $GOPATH/github.com/kubemq/demo/backend/elastic/kubemq-elastic-run .
COPY --from=builder $GOPATH/github.com/kubemq/demo/backend/elastic/.config.yaml .
EXPOSE 50000
CMD ["./kubemq-elastic-run"]

