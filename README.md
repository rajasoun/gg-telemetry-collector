# Git Guardian Secrets Telemetry Collector Scripts

## Getting Started 

1. Open Terminal In Dev Container

```
brew install jq dasel
cp .env.sample .env
touch cookie.txt
```

2. Edir `.env` and provide details for  GITGUARDIAN_URL and GITGUARDIAN_API_KEY. For
GITGUARDIAN_URL, provide just the domain without the protocol

3. Fetch Cookies from the browser and add to the `cookie.txt`

3. Run `collect.sh`