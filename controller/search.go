package controller

import (
	"github.com/labstack/echo"

	"net/http"
	"github.com/radepal/RPRS/yum"
	"fmt"
	"os"

)

func Search(c echo.Context) error {

	repo := yum.NewRepo()
	baseurl := "http://mirror.centos.org/centos/7/os/x86_64/"

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

	fmt.Printf("%v",repomd.Databases)

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
