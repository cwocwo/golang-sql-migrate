# golang-sql-migrate
## build
```
packr2 build main.go
```

```
packr2 build -o liqui -ldflags "-X main.version=1.5 -X main.gitcommit=112313213213213 -X main.buildstamp=20181214_152656" main/main.go
```