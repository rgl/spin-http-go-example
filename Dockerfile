# syntax=docker/dockerfile:1.6
FROM scratch
COPY ./dist/spin.toml .
COPY ./dist/spin-http-go-example.wasm .
