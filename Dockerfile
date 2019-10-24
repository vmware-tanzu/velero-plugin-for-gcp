# Copyright 2017, 2019 the Velero contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.13-buster AS build
WORKDIR /go/src/github.com/vmware-tanzu/velero-plugin-for-gcp
# copy vendor in separately so the layer can be cached if the contents don't change
COPY vendor vendor
COPY velero-plugin-for-gcp velero-plugin-for-gcp
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /go/bin/velero-plugin-for-gcp ./velero-plugin-for-gcp


FROM ubuntu:bionic
RUN mkdir /plugins
COPY --from=build /go/bin/velero-plugin-for-gcp /plugins/
USER nobody:nobody
ENTRYPOINT ["/bin/bash", "-c", "cp /plugins/* /target/."]