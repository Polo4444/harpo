package archiving

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestTar(t *testing.T) {

	// Init
	TestInit(t)

	// We create tar provider
	p, err := GetProvider(TarProvider, testTarConf)
	if err != nil {
		t.Fatalf("Error creating Zip provider: %s", err.Error())
	}

	// We create a zip file
	f, err := os.Create("test_dummies.tar.gz")
	if err != nil {
		t.Fatalf("Error creating file: %s", err.Error())
	}
	defer f.Close()

	// Create a context with a timeout of 1 minute
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// We archive the folder test_dummies
	err = p.Archive(ctx, "test_dummies", f, true)
	if err != nil {
		t.Fatalf("Error archiving: %s", err.Error())
	}

	TestEnd(t)
}
