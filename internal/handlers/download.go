package handlers

import (
	"errors"
	"fmt"
	"github.com/seanime-app/seanime/internal/anilist"
	"github.com/seanime-app/seanime/internal/downloader"
	"github.com/seanime-app/seanime/internal/nyaa"
	"github.com/seanime-app/seanime/internal/updater"
	"github.com/seanime-app/seanime/internal/util"
	"github.com/sourcegraph/conc/pool"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// HandleDownloadNyaaTorrents will get magnets from Nyaa and add them to qBittorrent.
// It also handles smart selection (downloader.SmartSelect).
//
//	POST /v1/download
func HandleDownloadNyaaTorrents(c *RouteCtx) error {

	type body struct {
		Urls        []string `json:"urls"`
		Destination string   `json:"destination"`
		SmartSelect struct {
			Enabled               bool  `json:"enabled"`
			MissingEpisodeNumbers []int `json:"missingEpisodeNumbers"`
			AbsoluteOffset        int   `json:"absoluteOffset"`
		} `json:"smartSelect"`
		Media *anilist.BaseMedia `json:"media"`
	}

	var b body
	if err := c.Fiber.BodyParser(&b); err != nil {
		return c.RespondWithError(err)
	}

	// try to start qbittorrent if it's not running
	err := c.App.QBittorrent.Start()
	if err != nil {
		return c.RespondWithError(err)
	}

	// get magnets
	p := pool.NewWithResults[string]().WithErrors()
	for _, url := range b.Urls {
		p.Go(func() (string, error) {
			return nyaa.TorrentMagnet(url)
		})
	}
	// if we couldn't get a magnet, return error
	magnets, err := p.Wait()
	if err != nil {
		return c.RespondWithError(err)
	}

	// create repository
	repo := &downloader.QbittorrentRepository{
		Logger:         c.App.Logger,
		Client:         c.App.QBittorrent,
		WSEventManager: c.App.WSEventManager,
		Destination:    b.Destination,
	}

	// try to add torrents to qbittorrent, on error return error
	err = repo.AddMagnets(magnets)
	if err != nil {
		return c.RespondWithError(err)
	}

	err = repo.SmartSelect(&downloader.SmartSelect{
		Magnets:               magnets,
		Enabled:               b.SmartSelect.Enabled,
		MissingEpisodeNumbers: b.SmartSelect.MissingEpisodeNumbers,
		AbsoluteOffset:        b.SmartSelect.AbsoluteOffset,
		Media:                 b.Media,
	})
	if err != nil {
		return c.RespondWithError(err)
	}

	return c.RespondWithData(true)

}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// HandleDownloadTorrentFile will download a torrent file from a given URL and save it to the destination folder.
//
//	POST /v1/download-torrent-file
func HandleDownloadTorrentFile(c *RouteCtx) error {

	type body struct {
		DownloadUrls []string           `json:"download_urls"`
		Destination  string             `json:"destination"`
		Media        *anilist.BaseMedia `json:"media"`
	}

	var b body
	if err := c.Fiber.BodyParser(&b); err != nil {
		return c.RespondWithError(err)
	}

	errs := make([]error, 0)
	for _, url := range b.DownloadUrls {
		err := downloadTorrentFile(url, b.Destination)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 1 {
		return c.RespondWithError(errs[0])
	} else if len(errs) > 1 {
		return c.RespondWithError(errors.New("failed to download multiple files"))
	}

	return c.RespondWithData(true)

}

func downloadTorrentFile(url string, dest string) (err error) {

	defer util.HandlePanicInModuleWithError("handlers/download/downloadTorrentFile", &err)

	// Get the file name from the URL
	fileName := filepath.Base(url)
	filePath := filepath.Join(dest, fileName)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful (status code 200)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, %s", resp.Status)
	}

	// Create the destination folder if it doesn't exist
	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// HandleDownloadRelease will download a release from a given URL and extract it to the destination folder.
//
//	POST /v1/download-release
func HandleDownloadRelease(c *RouteCtx) error {

	type retData struct {
		Destination string `json:"destination"`
		Error       string `json:"error,omitempty"`
	}

	type body struct {
		DownloadUrl string `json:"download_url"`
		Destination string `json:"destination"`
	}

	var b body
	if err := c.Fiber.BodyParser(&b); err != nil {
		return c.RespondWithError(err)
	}

	path, err := c.App.Updater.DownloadLatestRelease(b.DownloadUrl, b.Destination)

	if err != nil {
		if errors.Is(err, updater.ErrExtractionFailed) {
			return c.RespondWithData(retData{Destination: path, Error: err.Error()})
		}
		return c.RespondWithError(err)
	}

	return c.RespondWithData(retData{Destination: path})

}
