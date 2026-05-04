FROM scratch

EXPOSE 80

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/testsite-overview /testsite-overview

ENTRYPOINT ["/testsite-overview"]
