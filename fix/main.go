package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
)

func main() {
	chartName := "prometheus-community/prometheus"
	chartVersion := "13.8.0"
	releaseName := "prometheus"
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go template(chartName, chartVersion, releaseName, &wg)
	}
	wg.Wait()
}

func template(chartName, chartVersion, releaseName string, wg *sync.WaitGroup) {
	defer wg.Done()
	var valueOpts = new(values.Options)
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "", log.Printf); err != nil {
		log.Fatal(err)
	}
	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ReleaseName = releaseName
	client.Replace = true
	client.ClientOnly = true
	client.APIVersions = chartutil.VersionSet([]string{})
	client.IncludeCRDs = false
	client.Version = chartVersion
	cp, err := client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		log.Fatal(err)
	}
	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		log.Fatal(err)
	}
	chartRequested, err := loader.Load(cp)
	if err != nil {
		log.Fatal(err)
	}
	if req := chartRequested.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
				}
				if err := man.Update(); err != nil {
					log.Fatal(err)
				}
				if chartRequested, err = loader.Load(cp); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}
	}
	client.Namespace = settings.Namespace()
	_, err = client.Run(chartRequested, vals)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("release generated successfully")
}
