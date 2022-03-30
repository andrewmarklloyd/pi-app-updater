#!/bin/bash

set -euox pipefail

# TODO: move this to install script?

workDir="/home/runner/work/pi-app-deployer/pi-app-deployer"
envFile="/usr/local/src/.pi-app-deployer-agent.env"

if [[ $(whoami) != "root" ]]; then
  echo "Script must be run as root"
  exit 1
fi

if [[ -z ${HEROKU_API_KEY} ]]; then
  echo "HEROKU_API_KEY env var not set, exiting now"
  exit 1
fi

rm -f ${envFile}
cat <<< "HEROKU_API_KEY=${HEROKU_API_KEY}" > ${envFile}

mv ${workDir}/pi-app-deployer-agent /usr/local/src/
/usr/local/src/pi-app-deployer-agent install --appUser runneradmin --repoName ${repo} --manifestName ${manifestName} --envVar MY_CONFIG=testing --logForwarding

grep "MY_CONFIG\=testing" /usr/local/src/.pi-test-amd64.env >/dev/null
diff test/test-int-appconfigs.yaml /usr/local/src/.pi-app-deployer.appconfigs.yaml

sleep 10
journalctl -u pi-app-deployer-agent.service
journalctl -u pi-test-amd64.service
systemctl is-active pi-app-deployer-agent.service
systemctl is-active pi-test-amd64.service


journalctl -u pi-test-amd64.service -f
# git push to pi-test, check for new commit
