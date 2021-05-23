FROM alpine

ARG GLIBC_MIRROR=https://github.com/sgerrand/alpine-pkg-glibc
ARG GLIBC_VER=2.33-r0
RUN  apk add --no-cache --virtual=.build-deps curl tzdata && \
 curl -Lo /etc/apk/keys/sgerrand.rsa.pub "https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub" && \
 curl -Lo /glibc.apk "${GLIBC_MIRROR}/releases/download/${GLIBC_VER}/glibc-${GLIBC_VER}.apk" && \
 curl -Lo /glibc-bin.apk "${GLIBC_MIRROR}/releases/download/${GLIBC_VER}/glibc-bin-${GLIBC_VER}.apk" \
 && apk add --no-cache --allow-untrusted \
   /glibc.apk \
   /glibc-bin.apk \
 && rm /glibc.apk \
 && rm /glibc-bin.apk \
 && mkdir /etc/ddns

WORKDIR /opt/ddns

COPY ./build/dist/ddns ./ddns

CMD ["./ddns", "-c", "/etc/ddns"]