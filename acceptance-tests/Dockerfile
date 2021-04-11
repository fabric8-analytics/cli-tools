FROM registry.access.redhat.com/ubi8/go-toolset

USER root
WORKDIR /tests
COPY . /tests
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega/...
RUN yum -y update
RUN yum install -y maven
RUN yum install -y java-1.8.0-openjdk
RUN mvn --version
RUN yum install -y python36; yum clean all && ln -s /usr/bin/python3 /usr/bin/python
RUN yum install -y python3-pip
RUN python3 -V
RUN python3 -m pip -V
RUN pip3 -V
RUN rpm -i https://github.com/fabric8-analytics/cli-tools/releases/download/v0.1.0/crda_0.1.0_Linux-64bit.rpm

ENTRYPOINT [ "/tests/test.sh" ]