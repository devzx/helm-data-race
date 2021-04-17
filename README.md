# helm-data-race


Example code to reproduce the data race.

`helm` is required to be installed as a prerequisite.

## One time setup
```sh
# Add the prometheus repository
$ helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# Update
$ helm repo update
```

## Run
```sh
$ GORACE=history_size=7 go run -race main.go
```

Your output should be similar to:
```
==================
WARNING: DATA RACE
Read at 0x000004123f90 by goroutine 20:
  helm.sh/helm/v3/pkg/action.(*Install).Run()
      /Users/eli/go/pkg/mod/helm.sh/helm/v3@v3.5.3/pkg/action/install.go:200 +0x249
  main.template()
      /Users/eli/git/helm-data-race/main.go:84 +0x67b

Previous write at 0x000004123f90 by goroutine 22:
  helm.sh/helm/v3/pkg/action.(*Install).Run()
      /Users/eli/go/pkg/mod/helm.sh/helm/v3@v3.5.3/pkg/action/install.go:200 +0x337
  main.template()
      /Users/eli/git/helm-data-race/main.go:84 +0x67b

Goroutine 20 (running) created at:
  main.main()
      /Users/eli/git/helm-data-race/main.go:25 +0xed

Goroutine 22 (running) created at:
  main.main()
      /Users/eli/git/helm-data-race/main.go:25 +0xed
==================
release generated successfully
release generated successfully
release generated successfully
release generated successfully
release generated successfully
Found 1 data race(s)
exit status 66
```