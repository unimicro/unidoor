Requirements:

Execute `generate-token <name>` to generate a token for a new person.

Start by installing it as a service in a system that has systemd
installed:

    systemctl enable /root/www/door.service

And start the service:

    systemctl start door.service

And check the logs that everything went ok:

    journalctl -u door.service # add -f to continuously print new entries
