#! /bin/bash
set -e

SERVER=${1:-zero}

CMD='tar xzf - -C /tmp/
 && echo ≫ Making directory
 && mkdir -p /srv/unidoor/
 && echo ≫ Copying service
 && cp /tmp/unidoor.service /srv/unidoor/
 && echo ≫ Enabling service
 && systemctl enable /srv/unidoor/unidoor.service
'

echo 'Running command on "'${SERVER}'":' $CMD
tar czf - unidoor.service |ssh $SERVER $CMD
