Requirements:

Execute `generate-cert` to generate a self signed cert to get https.

Execute `generate-token <name>` to generate a token for a new person.

Start by symlinking the supervisor/door.conf script to supervisor and starting it there:

    ln -s /root/www/supervisor-door.conf /etc/supervisor/conf.d/door.conf

