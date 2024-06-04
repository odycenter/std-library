package ocr_test

import (
	"context"
	"fmt"
	"os"
	"std-library/sdk/google/ocr"
	"testing"
	"time"
)

func TestText(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	cert, _ := os.ReadFile("cert.json")
	f1, _ := os.Open("img.png")
	defer f1.Close()
	cli := ocr.NewImageAnnotatorClient(ctx, cert)
	annotations, err := cli.DetectTexts(f1, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, annotation := range annotations {
		fmt.Printf("%s\n", annotation.String())
	}
	f2, _ := os.Open("img.png")
	defer f2.Close()
	document, err := cli.DetectDocumentText(f2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", document.GetText())
}
