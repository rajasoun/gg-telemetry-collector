#!//usr/bin/env bash

# Purpose: Get Secrets for Repos from GG
# pre-requsites:
#   brew install jq dasel
# ------------------------------------------

source .env 
curl -s -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/sources?type=github&&page=1&per_page=200" | jq  > "data/data.json"
dasel -r json -w csv < "data/data.json" > "data/data.csv"

INPUT="data/data.csv"
OLDIFS=$IFS
IFS=','
[ ! -f $INPUT ] && { echo "$INPUT file not found"; exit 99; }

echo "id,full_name,secrets_count" > "result/data.csv"
sed 1d $INPUT | while IFS=, read -r full_name id	type url visibility
do
	echo "Getting Secrets For Repo $full_name with ID : $id"
    #secrets_count=$(curl -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/occurrences/secrets?source_id=$id"  |  jq -c '.[] | select(  .presence == "present")' | wc -l)
    secrets_count=$(curl -s -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/occurrences/secrets?source_id=$id"  |  jq -c  'map(.sha) | unique | length')
    echo "Secrets Count: $secrets_count"
    echo "$id,$full_name,$secrets_count" >> "result/data.csv"
done 
IFS=$OLDIFS



