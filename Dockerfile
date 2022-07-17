FROM golang:1.17-bullseye

FROM python:3.9-slim-bullseye
COPY --from=0 /usr/local/go/bin/go /usr/local/bin/go
COPY --from=0 /usr/local/go/bin/gofmt /usr/local/bin/gofmt
RUN mkdir /usr/local/go
COPY . /jroh
RUN pip install /jroh && rm --recursive --force /jroh "$(pip cache dir)"
ENTRYPOINT ["jrohc"]
CMD ["--help"]
