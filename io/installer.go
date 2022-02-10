package io

import (
	"fmt"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Resource struct {
	Filename string
	Url      string
}

func createProgressBar(fileSize int, progress *mpb.Progress) *mpb.Bar {
	bar := progress.AddBar(
		int64(fileSize),
		// Modification before progress bar
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f"), // Number of Downloads
			decor.Percentage(decor.WCSyncSpace),     // Progress percentage
		),
		// Modification after progress bar
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_GO, 90),
			decor.Name(" ] "),
			decor.EwmaSpeed(decor.UnitKiB, "% .2f", 60),
		),
	)
	return bar
}

func DownloadFile(resource Resource) (err error) {
	// Create the file
	target, err := os.Create(resource.Filename + ".tmp")
	finalPath := "./" + resource.Filename

	if err != nil {
		return err
	}
	defer func() {
		err := target.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Get the data
	resp, err := http.Get(resource.Url)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Check filesize
	fileSize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	progress := mpb.New(mpb.WithWidth(64), mpb.WithRefreshRate(180*time.Millisecond))
	bar := createProgressBar(fileSize, progress)
	reader := bar.ProxyReader(resp.Body)
	defer func() {
		err := reader.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Writer the body to temp file
	_, err = io.Copy(target, reader)
	if err != nil {
		return err
	}

	// Rename temp file to final file
	err = os.Rename(finalPath+".tmp", finalPath)
	if err != nil {
		return err
	}

	progress.Wait()
	return nil
}
