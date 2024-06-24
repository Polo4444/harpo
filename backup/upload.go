package backup

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Polo44444/harpo/alerting"
	"github.com/Polo44444/harpo/config"
	"github.com/Polo44444/harpo/storing"
)

var (
// uploadTimeout = time.Duration(1 * time.Hour) // TODO: Calculate the timeout based on the folder size
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

func (u *uploader) process(ctx context.Context, folder config.Folder, storages map[string]storing.Provider, notifiers map[string]alerting.Provider) {

	archiveFile, ok := ctx.Value(ArchiveCtxKey).(*os.File)

	if !ok {
		log.Printf("Unable to get archive file from context\n")
		NotifyError(
			ctx,
			folder.Name,
			fmt.Sprintf("Unable to get archive file of folder %s from context", folder.Name),
			"",
			nil,
			notifiers,
		)
		return
	}

	archiveFile.Close()
}
