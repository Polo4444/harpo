package backup

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Polo44444/harpo/alerting"
	"github.com/Polo44444/harpo/config"
	"github.com/Polo44444/harpo/storing"
)

var (
	uploadTimeout = time.Duration(90 * time.Minute) // TODO: Calculate the timeout based on the folder size
)

type uploader struct {
	next processor
}

func NewUploader() *uploader {
	return &uploader{}
}

func (u *uploader) setNext(p processor) processor {
	u.next = p
	return p
}

func (u *uploader) upload(
	ctx context.Context,
	wg *sync.WaitGroup,
	archiveFile string,
	contentType string,
	folder *config.Folder,
	notifiers map[string]alerting.Provider,
	storName string,
	stor storing.Provider) {

	defer wg.Done()

	// open file
	file, err := os.Open(archiveFile)
	if err != nil {
		log.Printf("Unable to open archive file %s: %v\n", archiveFile, err)
		NotifyError(
			ctx,
			folder.Name,
			fmt.Sprintf("Unable to open archive file %s", archiveFile),
			"",
			err,
			notifiers,
		)
		return
	}
	defer func() {
		file.Close()
	}()

	// Create timeout context
	uCtx, cancel := context.WithTimeout(ctx, uploadTimeout) // TODO: Calculate the timeout based on the folder size
	defer cancel()

	// Upload the archive file
	destFilePath := getDestFilePath(
		folder.Name,
		folder.Destination,
		filepath.Ext(file.Name()),
	)

	err = stor.UploadWithReader(uCtx, destFilePath, file, contentType)
	if err != nil {
		log.Printf("Unable to upload archive to storage %s: %v\n", storName, err)
		NotifyError(
			ctx,
			folder.Name,
			fmt.Sprintf("Unable to upload archive to storage %s", storName),
			"",
			err,
			notifiers,
		)
		return
	}

	NotifyInfo(
		ctx,
		folder.Name,
		fmt.Sprintf("Archive uploaded 📤 ✅ to storage %s", storName),
		"",
		notifiers,
	)
}

func (u *uploader) process(ctx context.Context, folder config.Folder, storages map[string]storing.Provider, notifiers map[string]alerting.Provider) {

	archiveFile, ok := ctx.Value(ArchiveCtxKey).(string)
	if !ok {
		log.Printf("Unable to get archive file from context\n")
		NotifyError(
			ctx,
			folder.Name,
			fmt.Sprintf("Unable to get archive file of folder 📁 %s from context", folder.Name),
			"",
			nil,
			notifiers,
		)
		return
	}

	// We get content type
	contentType, ok := ctx.Value(ContentTypeCtxKey).(string)
	if !ok {
		contentType = "application/octet-stream"
	}

	// ─── Start Upload Process ────────────────────────────────────────────
	NotifyInfo(
		ctx,
		folder.Name,
		fmt.Sprintf("Started archive upload 📤 of folder 📁 %s", folder.Name),
		"",
		notifiers,
	)

	// We will upload teh archive to all the storage providers.
	// Therefore, we will create multiple readers for each storage provider.
	// We will use the same archive file for each storage provider.

	var wg sync.WaitGroup
	for name, storage := range storages {

		wg.Add(1)
		go u.upload(ctx, &wg, archiveFile, contentType, &folder, notifiers, name, storage)
	}

	wg.Wait()

	if u.next != nil {
		u.next.process(ctx, folder, storages, notifiers)
	}
}
