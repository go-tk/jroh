FROM golang:1.16-alpine

FROM python:3.9.7-slim-bullseye
COPY --from=0 /usr/local/go/bin/gofmt /usr/local/bin/gofmt
COPY src /jroh/src
COPY setup.py /jroh/setup.py
RUN pip install --index-url https://mirrors.ustc.edu.cn/pypi/web/simple --no-cache-dir /jroh \
&& rm --recursive --force /root/src
ENTRYPOINT ["jrohc"]
CMD ["--help"]
