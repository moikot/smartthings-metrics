FROM --platform=$BUILDPLATFORM golang:1.14 as build-env

# xx wraps go to automatically configure $GOOS, $GOARCH, and $GOARM
# based on TARGETPLATFORM provided by Docker.
COPY --from=tonistiigi/xx:golang / /

ARG APP_FOLDER

ADD . ${APP_FOLDER}
WORKDIR ${APP_FOLDER}

# Compile independent executable using go wrapper from xx:golang
ARG TARGETPLATFORM
RUN CGO_ENABLED=0 go build -a -o /bin/main .

FROM scratch

ARG APP_FOLDER

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /bin/main /
COPY --from=build-env ${APP_FOLDER}/units.json /

ENTRYPOINT ["/main"]