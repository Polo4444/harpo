package backup

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Polo44444/harpo/alerting"
	"github.com/Polo44444/harpo/config"
	"github.com/Polo44444/harpo/storing"
)

type cleaner struct {
	next processor
}

func NewCleaner() *cleaner {
	return &cleaner{}
}

func (c *cleaner) setNext(p processor) processor {
	c.next = p
	return p
}

func (c *cleaner) process(ctx context.Context, folder config.Folder, storages map[string]storing.Provider, notifiers map[string]alerting.Provider) {

	defer func() {

		// We remove archive file
		archiveFile, ok := ctx.Value(ArchiveCtxKey).(string)
		if !ok {
			log.Printf("Unable to get archive file from context\n")
		}

		err := os.Remove(archiveFile)
		if err != nil {
			log.Printf("Unable to remove archive file %s: %v\n", archiveFile, err)
		}
	}()

	if !folder.Remove {
		c.success(ctx, folder, notifiers)
		return
	}

	// We get the folder path
	folderPath := folder.Path

	// We clear the folder content without removing the folder itself
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		log.Printf("Unable to read folder 📁 %s: %v\n", folderPath, err)
		NotifyError(
			ctx,
			folder.Name,
			fmt.Sprintf("Unable to read folder 📁 %s", folder.Name),
			"",
			err,
			notifiers,
		)
		return
	}

	for _, entry := range entries {

		fullPath := filepath.Join(folderPath, entry.Name())
		os.RemoveAll(fullPath) // No need to check errors here. It is a fault tolerant operation
	}

	log.Printf("Folder 📁 %s has been removed\n", folderPath)

	c.success(ctx, folder, notifiers)
}

func (c *cleaner) success(ctx context.Context, folder config.Folder, notifiers map[string]alerting.Provider) {

	NotifySuccess(
		ctx,
		folder.Name,
		fmt.Sprintf("Folder 📁 %s has been successfully backed up 💾 ✅ 🚀 🎉", folder.Name),
		"",
		notifiers,
	)
}
