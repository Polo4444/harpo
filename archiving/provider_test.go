package archiving

import (
	"archive/zip"
	"os"
	"testing"
)

var (
	testZipConf = BuildZipConfig(
		zip.Deflate,
	)
	testTarConf = BuildTarConfig(
		9,
		GzCompressionType,
	)
)

func TestInit(t *testing.T) {

	// We create a folder test_dummies, containing 2 files and a folder texts containing 1 file
	// We create file test_dummies/texts/file1.txt
	err := os.MkdirAll("test_dummies/texts", os.ModePerm)
	if err != nil {
		t.Fatalf("Error creating directory: %s", err.Error())
	}
	err = os.WriteFile("test_dummies/texts/file1.txt", []byte("Hello World!"), os.ModePerm)
	if err != nil {
		t.Fatalf("Error creating file: %s", err.Error())
	}

	// We create file test_dummies/file1.txt
	err = os.WriteFile("test_dummies/file1.txt", []byte("Hello World! 1"), os.ModePerm)
	if err != nil {
		t.Fatalf("Error creating file: %s", err.Error())
	}

	// We create file test_dummies/file2.txt
	err = os.WriteFile("test_dummies/file2.txt", []byte("Hello World! 2"), os.ModePerm)
	if err != nil {
		t.Fatalf("Error creating file: %s", err.Error())
	}
}

func TestEnd(t *testing.T) {

	// We remove the folder test_dummies
	err := os.RemoveAll("test_dummies")
	if err != nil {
		t.Fatalf("Error removing directory: %s", err.Error())
	}
}
