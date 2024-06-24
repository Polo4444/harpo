package backup

import (
	"archive/zip"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Polo44444/harpo/alerting"
	"github.com/Polo44444/harpo/archiving"
	"github.com/Polo44444/harpo/config"
	"github.com/Polo44444/harpo/storing"
	"github.com/google/uuid"
)

var (
	archiveTimeout = time.Duration(30 * time.Minute) // TODO: Calculate the timeout based on the folder size
)

type archiver struct {
	next processor
}

func NewArchiver() *archiver {
	return &archiver{}
}

func (a *archiver) setNext(p processor) processor {
	a.next = p
	return p
}

func (a *archiver) process(ctx context.Context, folder config.Folder, storages map[string]storing.Provider, notifiers map[string]alerting.Provider) {

	archiveType := strings.ToUpper(folder.Archiver)
	var (
		p   archiving.Provider
		err error
	)

	switch archiveType {
	case string(archiving.ZipProvider):
		p, err = archiving.GetProvider(archiving.ZipProvider, archiving.BuildZipConfig(zip.Deflate))
	case string(archiving.TarProvider):
		p, err = archiving.GetProvider(archiving.TarProvider, archiving.BuildTarConfig(9, archiving.GzCompressionType))
	default:
		p, err = archiving.GetProvider(archiving.ZipProvider, archiving.BuildZipConfig(zip.Deflate))
	}

	if err != nil {
		log.Printf("Unable to get archiver provider of folder %s: %v\n", err, folder.Path)
		NotifyError(
			ctx,
			folder.Name,
			fmt.Sprintf("Unable to get archiver provider of folder %s", folder.Name),
			"",
			err,
			notifiers,
		)
		return
	}

	// ─── Start Archiving Process ─────────────────────────────────────────
	NotifyInfo(
		ctx,
		folder.Name,
		fmt.Sprintf("Backup 💾 process of folder %s started🌴", folder.Name),
		"",
		notifiers,
	)

	// Create a file to write the archive
	fileName := uuid.Must(uuid.NewRandom()).String() + p.Ext()
	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Unable to create archive file of folder %s: %v\n", err, folder.Path)
		NotifyError(
			ctx,
			folder.Name,
			fmt.Sprintf("Unable to create archive file of folder %s", folder.Name),
			"",
			err,
			notifiers,
		)
		return
	}

	pCtx, cancel := context.WithTimeout(ctx, archiveTimeout) // TODO: Calculate the timeout based on the folder size
	defer cancel()
	err = p.Archive(pCtx, folder.Path, file, folder.IgnoreArchiveErrors)
	if err != nil {
		log.Printf("Unable to archive folder %s: %v\n", err, folder.Path)
		NotifyError(
			ctx,
			folder.Name,
			fmt.Sprintf("Unable to archive folder %s", folder.Name),
			"",
			err,
			notifiers,
		)
		return
	}

	// ─── End Archiving Process ───────────────────────────────────────────
	NotifyInfo(
		ctx,
		folder.Name,
		fmt.Sprintf("Archival of folder %s completed🗜️✅", folder.Name),
		"",
		notifiers,
	)

	// Hold the file inside the context
	newCtx := context.WithValue(ctx, ArchiveCtxKey, file)

	if a.next != nil {
		a.next.process(newCtx, folder, storages, notifiers)
	}
}
