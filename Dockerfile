FROM alpine:3.16

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && apk add -U tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN mkdir /app
ENV PATH="/app:${PATH}"
ENV TZ=Asia/Shanghai
WORKDIR /app
ADD boltdb-explorer /app/boltdb-explorer
EXPOSE 8080
CMD [ "boltdb-explorer" ]