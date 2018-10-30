#! /bin/bash
set -e

SERVER=${1:-zero}

CMD='
tar xzf - -C /srv/unidoor/
 && echo ≫ Replacing executable...
 && test -f /srv/unidoor/unidoor
 && mv /srv/unidoor/unidoor{,.old} || true
 && mv /srv/unidoor/unidoor{.new,}
 && echo ≫ Restarting service...
 && systemctl daemon-reload
 && service unidoor restart
 && echo ≫ Checking status...
 && service unidoor status
 && echo ≫ Done
'

echo '≫ Transmitting executable "unidoor.new" to "'${SERVER}'" remote...'
tar czf - unidoor.new unidoor.service add-token | ssh $SERVER $CMD
