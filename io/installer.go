package io

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"svm/parsers"
	"time"
)

type Resource struct {
	Home    string
	Url     string
	Version parsers.Version
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
	tempFilePath := filepath.Join(resource.Home, ".svm", resource.Version.SparkVersion+".tmp")
	tempFile, err := os.Create(tempFilePath)
	tarPath := filepath.Join(resource.Home, ".svm", resource.Version.SparkVersion+".tgz")
	if err != nil {
		return err
	}
	defer func() {
		err := tempFile.Close()
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
	_, err = io.Copy(tempFile, reader)
	if err != nil {
		return err
	}

	// Rename temp file to final file
	err = os.Rename(tempFilePath, tarPath)
	if err != nil {
		return err
	}

	progress.Wait()
	return nil
}

// UnzipTar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func UnzipTar(resource Resource) error {
	svmPath := filepath.Join(resource.Home, ".svm")
	tarPath := filepath.Join(resource.Home, ".svm", resource.Version.SparkVersion+".tgz")
	tarFile, err := os.Open(tarPath)
	gzr, err := gzip.NewReader(tarFile)
	if err != nil {
		return err
	}
	defer func(gzr *gzip.Reader) {
		err := gzr.Close()
		err = os.Remove(tarPath)
		if err != nil {

		}
	}(gzr)

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(svmPath, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			err = f.Close()
			if err != nil {
				return err
			}
		}
	}

}

func RenameUnzipped(resource Resource) error {
	sparkUnzippedPath := filepath.Join(resource.Home, ".svm", resource.Version.FullVersion)
	sparkFinalPath := filepath.Join(resource.Home, ".svm", resource.Version.SparkVersion)
	err := os.Rename(sparkUnzippedPath, sparkFinalPath)
	if err != nil {
		return err
	}
	return nil
}
