FROM golang:1.6-onbuild

EXPOSE 8080

ENV SPOT_FLEET_REQUEST_IDS "sfr-xxxxx-xxx-xxx-xxxx,sfr-yyyyy-yyy-yy-yyyy"
ENV TAGS "Owner:user@domain,Cluster:us-west-2"

ENV KEEP_RUNNING "true"
ENV SLEEP_INTERVAL "60"
