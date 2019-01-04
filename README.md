# Install
```bash
git clone git@github.com:axklim/envci.git .
go build -o /usr/local/bin/envci ./src
envci
the required flags `-p, --project', `-t, --token' and `-u, --url' were not specified
```

# Use
List user-defined GitLab CI/CD Variables
```bash
envci -p gudik/envci-demo -u https://gitlab.com/api/v4 -t <GITLAB_TOKEN>
TEST_GITLAB_VARIABLE='some values'
MODE='debug'
``` 

Set variables as environment current console session
```bash
. <(envci -p gudik/envci-demo -u https://gitlab.com/api/v4 -t <GITLAB_TOKEN>)
```

Clear environments
```bash
. <(envci -p gudik/envci-demo -u https://gitlab.com/api/v4 -t <GITLAB_TOKEN> --clear)
```
