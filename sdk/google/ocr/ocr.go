package ocr

import (
	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"context"
	"google.golang.org/api/option"
	"io"
)

type ImageAnnotatorClient struct {
	ctx context.Context
	cli *vision.ImageAnnotatorClient
	err error
}

func NewImageAnnotatorClient(ctx context.Context, cert []byte) *ImageAnnotatorClient {
	client, err := vision.NewImageAnnotatorClient(ctx,
		option.WithCredentialsJSON(cert))
	return &ImageAnnotatorClient{ctx, client, err}
}

// DetectTexts 对图像执行文本检测。最多返回 maxResults 个结果。
func (i *ImageAnnotatorClient) DetectTexts(img io.Reader, maxResults int) ([]*visionpb.EntityAnnotation, error) {
	if i.err != nil {
		return nil, i.err
	}
	image, err := vision.NewImageFromReader(img)
	if err != nil {
		return nil, err
	}

	return i.cli.DetectTexts(i.ctx, image, nil, maxResults)
}

// DetectDocumentText 对图像执行全文 (OCR) 检测。
func (i *ImageAnnotatorClient) DetectDocumentText(img io.Reader) (*visionpb.TextAnnotation, error) {
	if i.err != nil {
		return nil, i.err
	}
	image, err := vision.NewImageFromReader(img)
	if err != nil {
		return nil, err
	}

	return i.cli.DetectDocumentText(i.ctx, image, nil)
}
