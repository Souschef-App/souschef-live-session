# Live Session

### GO Commands

#### Build

```shell
go build [filename]
```

#### Run

```shell
go run [filename]
```

### Docker Commands

#### Build

```shell
# -t: Image tag, in our case "souschef-live-session"
docker build -t [image_name] .
```

#### Run

```shell
# Note: container will generate random name

# -d: Run container in the background
# --rm: Auto-delete when stopped
docker run -d --rm -p 8080:8080 souschef-live-session
```
