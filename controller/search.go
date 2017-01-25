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
	requestedpackage := c.Param("package")
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
	repo.ID = requestedrepo
	// get cache for this repo
	_, err := repo.CacheLocal("cache")
	if err != nil {
		return err
	}

	file, err := os.Open(fmt.Sprintf("cache/%s/gen/primary.xml", requestedrepo))
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

	var result []string
	packages := primarymeta.Packages
	for _, element := range packages {
		if element.Name() == requestedpackage {
			result = append(result, element.String())
		}
	}
	type Response struct {
		Result []string `json:"result"`
	}
	response := Response{Result: result}
	return c.JSON(http.StatusOK, response)
}
