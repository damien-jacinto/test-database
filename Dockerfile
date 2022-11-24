FROM golang:1.18-alpine as build
ARG APPUSER=appuser

ENV USER=${APPUSER}
ENV UID=1001

RUN adduser -D -g "" -H -s "/sbin/nologin" -u "${UID}" "${USER}"
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM scratch
ARG APPUSER=appuser

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER ${APPUSER}:${APPUSER}
COPY --from=build --chown=${APPUSER}:${APPUSER} /app/main /go/bin/

WORKDIR /home/${APPUSER}

EXPOSE 8080

ENTRYPOINT ["/go/bin/main"]
