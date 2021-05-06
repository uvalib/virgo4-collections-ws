#!/usr/bin/env bash
#
# run application
#

# run the server
cd bin; ./virgo4-collections-ws -solr $COLLECTIONS_SOLR_URL -core $COLLECTIONS_CORE_NAME

# return the status
exit $?

#
# end of file
#
