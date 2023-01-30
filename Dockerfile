FROM golang:1.19 AS base
ENTRYPOINT [ "/main" ]
COPY /main /
