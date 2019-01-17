# Setup:

Setup ssh autologin with the name "unidoor" in ~/.ssh/config

Run "./setup.sh" to setup the folders and service on a new server.

Run "./build.sh && ./deploy.sh" every time you want to deploy the service.

If you have another name for the server in ssh/config you can put that after the deploy script,
i.e. `./build.sh && ./deploy.sh <my server name>`

PS: This project uses git submodules for it's go dependencies.

# Usage:

Execute `add-token <name>` to generate a token for a new person.

Check the logs that works ok:

    journalctl -u unidoor # add -f to continuously print new entries
