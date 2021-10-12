#
# build the target application
#
FROM golang:1.16.5-alpine3.13 as builder
RUN apk update && apk upgrade && apk add make yarn

WORKDIR /build
COPY go.mod go.sum Makefile ./
COPY frontend ./frontend
COPY backend ./backend
RUN make linux-full

#
# build the target container
#
FROM alpine:3.13

# update the packages
RUN apk update && apk upgrade && apk add bash tzdata ca-certificates curl && rm -fr /var/cache/apk/*

# Create the run user and group
RUN addgroup webservice && adduser webservice -G webservice -D

# set the timezone appropriatly
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Specify home
ENV APP_HOME /virgo4-collections-ws
WORKDIR $APP_HOME

# Create necessary directories
RUN mkdir -p $APP_HOME/scripts $APP_HOME/bin
RUN chown -R webservice $APP_HOME && chgrp -R webservice $APP_HOME

# add the migrate tool
RUN cd $APP_HOME/bin && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz && mv migrate.linux-amd64 migrate

# port and run command
EXPOSE 8080
CMD scripts/entry.sh

# Move in necessary assets
COPY package/scripts/entry.sh $APP_HOME/scripts/entry.sh
COPY package/scripts/migrate.sh $APP_HOME/scripts/migrate.sh
COPY db/migrations/*.sql $APP_HOME/db/
COPY package/data/container_bash_profile /home/docker/.profile
COPY --from=builder /build/bin/public $APP_HOME/bin/public
COPY --from=builder /build/bin/virgo4-collections-ws.linux $APP_HOME/bin/virgo4-collections-ws

# Ensure permissions are correct
RUN chown webservice:webservice /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/scripts/migrate.sh  $APP_HOME/db $APP_HOME/bin/virgo4-collections-ws && chmod 755 /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/bin/virgo4-collections-ws $APP_HOME/scripts/migrate.sh  $APP_HOME/db
RUN chown -R webservice:webservice $APP_HOME/bin/public && chmod -R 755 $APP_HOME/bin/public


# Add the build tag
ARG BUILD_TAG
RUN test -n "$BUILD_TAG" && touch $APP_HOME/buildtag.build-$BUILD_TAG || touch $APP_HOME/buildtag.build-0

# Specify the user
USER docker

#
# end of file
#
