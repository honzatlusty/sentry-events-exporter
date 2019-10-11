#!/bin/bash

token=$1
sentry_url=$2
hours_ago=$3
timestamp=$(date +%s --date=" ${hours_ago} hours ago")
timestamp_now=$(date +%s)
organizations=$(curl --silent -H "Authorization: Bearer ${token}" ${sentry_url}/api/0/organizations/ | jq .[].slug -r)

echo "# HELP sentry_events_received_per_hour_total /api/0/projects/<organization>/<project>/stats/?&since=${timestamp}&until=${timestamp_now}
# TYPE sentry_events_received_per_hour_total gauge"

for organization in $organizations; do
  projects=$(curl --silent -H "Authorization: Bearer ${token}" "${sentry_url}/api/0/organizations/${organization}/projects/" | jq .[].slug -r)
  for project in $projects; do
    events=$(curl --silent -H "Authorization: Bearer ${token}" "${sentry_url}/api/0/projects/${organization}/${project}/stats/?&since=${timestamp}&until=${timestamp_now}" | jq .[][1] | awk '{s+=$1} END {print s}')
    project_sanitized=$(echo $project | tr ' -.' '_')
    echo "sentry_events_received_per_hour_total{project=\"${project_sanitized}\"} ${events}"
  done
done
