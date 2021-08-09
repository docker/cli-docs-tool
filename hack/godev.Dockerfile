# syntax=docker/dockerfile:1.3

# Copyright 2021 docgen authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG GO_VERSION
ARG GITHUB_REPOSITORY
ARG GITHUB_REF

FROM golang:${GO_VERSION}-alpine AS godev
ARG GITHUB_REPOSITORY
ARG GITHUB_REF
ENV GOPROXY=https://proxy.golang.org
RUN if [ -z "${GITHUB_REPOSITORY}" -o -z "${GITHUB_REF}" ]; then echo >&2 "GITHUB_REPOSITORY and GITHUB_REF required"; exit 1; fi; \
  go get github.com/${GITHUB_REPOSITORY}@${GITHUB_REF#refs/tags/}
