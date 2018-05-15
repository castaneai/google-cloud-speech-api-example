# google-cloud-speech-api-example

## Streaming example

```bash
ffmpeg -re -i <input> -f s16le -ar 48k -ac 1 -loglevel 8 pipe:1 | go run streaming/main.go
```