#
# build the target application
#
FROM public.ecr.aws/docker/library/golang:1.23.4-alpine3.21 AS builder

# update the packages
RUN apk update && apk upgrade && apk add --no-cache make npm yarn

WORKDIR /build
COPY go.mod go.sum Makefile ./
COPY frontend ./frontend
COPY backend ./backend
RUN make linux-full

#
# build the target container
#
FROM public.ecr.aws/docker/library/alpine:3.21

# update the packages
RUN apk update && apk upgrade && apk add bash tzdata ca-certificates curl exiftool && rm -fr /var/cache/apk/*

# Create the run user and group
RUN addgroup webservice && adduser webservice -G webservice -D

# set the timezone appropriatly
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Specify home
ENV APP_HOME=/virgo4-collections-ws
WORKDIR $APP_HOME

# Create necessary directories
RUN mkdir -p $APP_HOME/scripts $APP_HOME/bin
RUN chown -R webservice $APP_HOME && chgrp -R webservice $APP_HOME

# add the migrate tool
ENV MIGRATE_VERSION=v4.18.1
RUN cd $APP_HOME/bin && curl -L https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz | tar xvz && rm LICENSE README.md

# port and run command
EXPOSE 8080
ENTRYPOINT ["scripts/entry.sh"]

# Move in necessary assets
COPY package/scripts/entry.sh $APP_HOME/scripts/entry.sh
COPY package/scripts/migrate.sh $APP_HOME/scripts/migrate.sh
COPY db/migrations/*.sql $APP_HOME/db/
COPY package/data/container_bash_profile /home/webservice/.profile
COPY --from=builder /build/bin/public $APP_HOME/bin/public
COPY --from=builder /build/bin/virgo4-collections-ws.linux $APP_HOME/bin/virgo4-collections-ws

# Ensure permissions are correct
RUN chown webservice:webservice /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/scripts/migrate.sh  $APP_HOME/db $APP_HOME/bin/virgo4-collections-ws && chmod 755 /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/bin/virgo4-collections-ws $APP_HOME/scripts/migrate.sh  $APP_HOME/db
RUN chown -R webservice:webservice $APP_HOME/bin/public && chmod -R 755 $APP_HOME/bin/public


# Add the build tag
ARG BUILD_TAG
RUN test -n "$BUILD_TAG" && touch $APP_HOME/buildtag.build-$BUILD_TAG || touch $APP_HOME/buildtag.build-0

# Specify the user
USER webservice

#
# end of file
#
