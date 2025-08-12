FROM docker.io/golang:1.24.6 AS build

RUN go env -w CGO_ENABLED=0

COPY src /app
WORKDIR /app
RUN go build -o /out/lo-test-task lo-test-task/cmd/server

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/lo-test-task /lo-test-task

CMD ["/lo-test-task"]