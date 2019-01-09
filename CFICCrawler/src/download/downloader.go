package download

import (
	"routingpool"
	"routingpool/task"
)

type DownloadInfo struct {
	Url string
	Path string
	Overwrite bool
}

type Downloader struct {
	httpRoutingPool *routingpool.ThreadPool
}

var downloader *Downloader

func (d *Downloader) Download(info *DownloadInfo) {
	caller := func(id int) {
		r := &Request{Url: info.Url, File: info.Path, OverWrite: info.Overwrite}
		_, err := r.Get()
		if err != nil {
			logger.Errorf("Request failure, %s", err)
			return
		}
	}

	d.httpRoutingPool.PutTask(task.NewDownloadTask("HTTP-Downloader", caller))
}

func NewDownloader(cfg *DownloadConfig) *Downloader {
	if downloader != nil {
		return downloader
	}

	downloader := &Downloader{}
	downloader.httpRoutingPool = routingpool.GetPool(128, 128)

	return downloader
}