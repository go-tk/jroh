FROM golang:1.16-alpine

FROM python:3.9.7-slim-bullseye
COPY --from=0 /usr/local/go/bin/gofmt /usr/local/bin/gofmt
COPY . /jroh
ARG VERSION
RUN ([ -z ${VERSION} ] || echo "VERSION = \"${VERSION}\"" >/jroh/src/jroh/version.py) \
&& pip install --no-cache-dir /jroh \
&& rm --recursive --force /jroh
ENTRYPOINT ["jrohc"]
CMD ["--help"]
