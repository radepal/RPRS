package controller

import (
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/radepal/go-yum"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

func Search(c echo.Context) error {

	requestedrepo := c.Param("repo")
	if !viper.IsSet("Repos") {
		return errors.New("Missing Repos config")
	}
	repos := viper.Sub("Repos")
	if !viper.IsSet("Repos") {
		return errors.New("Missing Repos config")
	}
	if !repos.IsSet(requestedrepo) {
		return errors.New(fmt.Sprintf("Missing repo %s config", requestedrepo))
	}
	repoconfig := repos.Sub(requestedrepo)

	repo := yum.NewRepo()
	baseurl := repoconfig.GetString("baseurl")

	repo.BaseURL = baseurl

	// get cache for this repo
	_, err := repo.CacheLocal("cache")
	if err != nil {
		return err
	}

	// repo metadata is always found at the repodata/repomd.xml subpath
	repomdurl := baseurl + "repodata/repomd.xml"

	// get repo metadata from url
	resp, err := http.Get(repomdurl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// decode http stream into repo metadata struct
	repomd, err := yum.ReadRepoMetadata(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", repomd.Databases)

	file, err := os.Open("cache/gen/primary.xml")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// decode http stream into repo metadata struct
	primarymeta, err := yum.ReadPrimaryMetadata(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Downloaded repository metadata revision %d\n", primarymeta.PackagesCount)

	return c.JSON(http.StatusOK, primarymeta)
}
