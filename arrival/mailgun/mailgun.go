package mailgun

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/kaz/rusuden/arrival"
	"github.com/mailgun/mailgun-go/v4"
)

type (
	MailgunArrival struct {
		client mailgun.Mailgun
	}
)

func New(key string) arrival.Arrival {
	return &MailgunArrival{mailgun.NewMailgun("", key)}
}

func (ma *MailgunArrival) Arrive(req *http.Request) ([]byte, error) {
	mediaType, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("mime.ParseMediaType failed: %w", err)
	}

	switch mediaType {
	case "multipart/form-data":
		return ma.arriveViaForward(req, params["boundary"])
	default:
		return nil, fmt.Errorf("Unexpected media type: %v", mediaType)
	}
}

func (ma *MailgunArrival) arriveViaForward(req *http.Request, boundary string) ([]byte, error) {
	mr := multipart.NewReader(req.Body, boundary)

	meta := map[string][]byte{}
	attachments := map[string][]byte{}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("mr.NextPart failed: %w", err)
		}

		data, err := ioutil.ReadAll(part)
		if err != nil {
			return nil, fmt.Errorf("ioutil.ReadAll failed: %w", err)
		}

		_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
		if err != nil {
			return nil, fmt.Errorf("mime.ParseMediaType failed: %w", err)
		}

		if _, ok := params["name"]; ok {
			meta[params["name"]] = data
		}
		if _, ok := params["filename"]; ok {
			attachments[params["filename"]] = data
		}
	}

	sig := mailgun.Signature{
		TimeStamp: string(meta["timestamp"]),
		Token:     string(meta["token"]),
		Signature: string(meta["signature"]),
	}
	verified, err := ma.client.VerifyWebhookSignature(sig)
	if err != nil {
		return nil, fmt.Errorf("ma.client.VerifyWebhookSignature failed: %w", err)
	}
	if !verified {
		return nil, fmt.Errorf("invalid signature")
	}

	for filename, data := range attachments {
		if strings.HasSuffix(filename, ".wav") {
			return data, nil
		}
	}
	return nil, fmt.Errorf("attachment not found")
}
