package storing

import (
	"context"
	"io"

	"github.com/Polo44444/harpo/models"
	storing_s3 "github.com/Polo44444/harpo/storing/s3"
)

const (
	S3Provider models.ProviderEntity = "S3"
)

// Provider interface
type Provider interface {

	/*Test tests the connection to the provider.
	 */
	Test(ctx context.Context) error

	/*UploadWithReader uploads the data to the filePath.
	`filePath` is the path where the data will be uploaded.
	`data` is the data to be uploaded.
	`contentType` is the content type of the data. Default is "application/octet-stream".
	*/
	UploadWithReader(ctx context.Context, filePath string, data io.Reader, contentType string) error

	/*DownloadWithWriter downloads the data from the filePath.
	`filePath` is the path where the data will be downloaded.
	`writer` is the writer where the data will be written.
	*/
	DownloadWithWriter(ctx context.Context, filePath string, writer io.Writer) error

	/*Info returns the information of the file in the filePath.
	`filePath` is the path where the data is located.
	*/
	Info(ctx context.Context, filePath string) (*models.FileInfo, error)

	/*Delete deletes the data from the filePath.
	`filePath` is the path where the data will be deleted.
	*/
	Delete(ctx context.Context, filePath string) error

	/*DeleteMany removes multiple files from the storage.
	`filePaths` is a list of paths where the data will be deleted.
	*/
	// If one file fails to be deleted, the function will continue deleting the rest of the files.
	// If all files fails to be deleted, the function will return an nil (no error).
	DeleteMany(ctx context.Context, filePaths []string) error

	/*Close closes the provider.
	 */
	Close(ctx context.Context) error
}

// GetProvider returns a provider based on the entity and the config
func GetProvider(entity models.ProviderEntity, config models.ProviderConfig) (Provider, error) {

	var err error = nil
	var prvd Provider = nil

	switch entity {
	case S3Provider:
		prvd, err = storing_s3.NewS3Provider(config)
	default:
		err = models.ErrProviderNotSupported
	}

	return prvd, err
}
