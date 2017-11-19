# start from hypriot/rpi-alpine-scratch (nginx:alpine)
#
# -------------------------------------------------
# FROM resin/raspbian
# MAINTAINER kjell.almgren[at]tetracon.se
# -------------------------------------------------
#
# FROM resin/rpi-raspbian
FROM alpine

MAINTAINER kjell.almgren@tetracon.se

# make some update to the OS in the container
#RUN apk update && \
#apk upgrade && \
#apk add bash && \
#rm -rf /var/cache/apk/*

#make some changes to the container images (docker dns-bugs)
#COPY docker-compose.yml docker-compose.yaml
#switch to our app directory (/selmametrics)
RUN mkdir -p /selmametrics
WORKDIR /selmametrics

# COPY executable selmametrics /selmametrics
COPY selmametrics /selmametrics

# copy our self-signed certificate
##COPY tetracon-server.crt /go/src/selmametrics
##COPY tetracon-server.key /go/src/selmametrics

# tell we are exposing our service on port 8000
EXPOSE 8000

# run it!

ENTRYPOINT ["./selmametrics"]
#CMD ["./selmametrics"]