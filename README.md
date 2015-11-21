# website-build
hugo website automation build aplication

Custom server that will wait for github push notification and then pull web site project and run hugo to rebuild content.

## Configuration

Set following Env variables:
- PUSH_PORT what port to lisen
- PUSH_HTTP_PATH http path to react to eq. /push
- PUSH_GITHUB_SECRET github webhook secret
- PUSH_PROJECT_DIR dir path where project is

## Dependencies

```bash
go get github.com/phayes/hookserve
```
