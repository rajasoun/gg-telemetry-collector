#!//usr/bin/env bash

# Purpose: Get Secrets for Repos from GG
# pre-requsites:
#   brew install jq dasel
# ------------------------------------------

source .env 
GITGUARDIAN_API_URL="https://$GITGUARDIAN_URL/exposed"

curl -s -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/sources?type=github&&page=1&per_page=100" | jq  > "data/data.json"
curl -s -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/sources?type=github&&page=2&per_page=100" | jq  >> "data/data.json"
curl -s -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/sources?type=github&&page=3&per_page=100" | jq  >> "data/data.json"
dasel -r json -w csv < "data/data.json" > "data/data.csv"

INPUT="data/data.csv"
OLDIFS=$IFS
IFS=','
[ ! -f $INPUT ] && { echo "$INPUT file not found"; exit 99; }

function get_secrets_count_from_api(){
    echo "id,full_name,secrets_count" > "result/data.csv"
    sed 1d $INPUT | while IFS=, read -r full_name id	type url visibility
    do
        echo "Getting Secrets For Repo $full_name with ID : $id"
        secrets_count=$(curl -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/occurrences/secrets?source_id=$id"  |  jq -c '.[] | select(  .presence == "present")' | wc -l)
        #secrets_count=$(curl -s -H "Authorization: Token ${GITGUARDIAN_API_KEY}" "${GITGUARDIAN_API_URL}/v1/occurrences/secrets?source_id=$id"  |  jq -c  'map(.sha) | unique | length')
        echo "Secrets Count: $secrets_count"
        echo "$id,$full_name,$secrets_count" >> "result/data.csv"
    done 
    IFS=$OLDIFS
}

function get_secrets_count_from_webcall(){
    local http_url="https://$GITGUARDIAN_URL/api/v1/accounts/2/sources/?monitored=true&page=1&page_size=10&search="
    local cookie=$(<cookie.txt) 

    echo "id,full_name,secrets_count" > "result/data.csv"
    sed 1d $INPUT | while IFS=, read -r full_name id	type url visibility
    do
        echo "Getting Secrets Count For Repo $full_name"
        filter_criteria=$(printf %s "$full_name" | jq -sRr @uri)
        web_end_point="$http_url$filter_criteria&ordering=-open_issues_count"
        json=$(curl -s $web_end_point -H $cookie --compressed)
        secrets_count=$(jq -r --arg repo_name $full_name '.results[] | select(.url | endswith($repo_name) ) | .open_issues_count' <<< "${json}" )
        if [ ! -z $secrets_count ] ;then 
            if [ "$secrets_count" -ne "0"  ] ;then
                echo "$id,$full_name,$secrets_count" >> "result/data.csv"
            fi 
        fi 
    done 
    IFS=$OLDIFS
}

get_secrets_count_from_webcall



