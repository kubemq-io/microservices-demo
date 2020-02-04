FROM kubemq/gobuilder as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq/demo/backend/users-api
WORKDIR $GOPATH/github.com/kubemq/demo/backend/users-api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -o kubemq-users-api-run .
FROM alpine:latest
RUN apk add ca-certificates
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
COPY --from=builder $GOPATH/github.com/kubemq/demo/backend/users-api/kubemq-users-api-run .
COPY --from=builder $GOPATH/github.com/kubemq/demo/backend/users-api/.config.yaml .
EXPOSE 50000
CMD ["./kubemq-users-api-run"]

