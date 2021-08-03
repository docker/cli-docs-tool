# syntax=docker/dockerfile:1.3
ARG GO_VERSION
ARG GITHUB_REPOSITORY
ARG GITHUB_REF

FROM golang:${GO_VERSION}-alpine AS godev
ARG GITHUB_REPOSITORY
ARG GITHUB_REF
ENV GOPROXY=https://proxy.golang.org
RUN if [ -z "${GITHUB_REPOSITORY}" -o -z "${GITHUB_REF}" ]; then echo >&2 "GITHUB_REPOSITORY and GITHUB_REF required"; exit 1; fi; \
  go get github.com/${GITHUB_REPOSITORY}@${GITHUB_REF#refs/tags/}
