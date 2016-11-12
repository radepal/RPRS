package controller

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Upload(c echo.Context) error {
	// Read form fields
	repo := c.FormValue("repo")

	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("data")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Crate directory
	path := filepath.Join(viper.GetString("UploadRpmPath"), repo)
	os.MkdirAll(path, os.ModeDir)

	// Destination
	dst, err := os.Create(path + string(os.PathSeparator) + file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("File %s uploaded successfully to repo %s", file.Filename, repo))
}
