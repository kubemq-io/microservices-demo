FROM kubemq/gobuilder as builder
ARG VERSION
ARG GIT_COMMIT
ARG BUILD_TIME
ENV GOPATH=/go
ENV PATH=$GOPATH:$PATH
ENV ADDR=0.0.0.0
ADD . $GOPATH/github.com/kubemq/demo/backend/users-web
WORKDIR $GOPATH/github.com/kubemq/demo/backend/users-web
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -mod=vendor -installsuffix cgo -o kubemq-users-web-run .
FROM alpine:latest
RUN apk add ca-certificates
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
COPY --from=builder $GOPATH/github.com/kubemq/demo/backend/users-web/kubemq-users-web-run .
COPY --from=builder $GOPATH/github.com/kubemq/demo/backend/users-web/.config.yaml .
COPY --from=builder $GOPATH/github.com/kubemq/demo/backend/users-web/web/users/dist/users  ./web/users/dist/users
CMD ["./kubemq-users-web-run"]

