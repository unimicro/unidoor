# Setup:

Setup ssh autologin with the name "zero" in ~/.ssh/config

Run "./setup.sh"

Run "./build.sh && ./deploy.sh"

If you have another name for the server in ssh/config you can put that after the deploy script,
i.e. `./build.sh && ./deploy.sh <my server name>`

PS: This project uses git submodules for it's go dependencies.

# Usage:

Execute `add-token <name>` to generate a token for a new person.

Check the logs that works ok:

    journalctl -u unidoor # add -f to continuously print new entries
