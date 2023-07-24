# Imgram

## Environment Variables
- `IMGRAM_DATA_DIR`: directory used to store uploaded images.
- `IMGRAM_DSN`: Database DSN used for connecting.

## Local setup

### Database setup
- Create user:password (defaults to root:)
- Create database (defaults to `imgram`)

### Database migration
```
migrate -path=migrations -database="$DSN" up
```

### Running
```
go run .
```

## Building
```
go build -o imgram .
```
