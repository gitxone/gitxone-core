FROM golang:latest

RUN useradd gitxone --create-home --shell /bin/bash
WORKDIR /home/gitxone

