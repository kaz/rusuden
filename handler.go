package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mgArr "github.com/kaz/rusuden/arrival/mailgun"
	mgDep "github.com/kaz/rusuden/departure/mailgun"
	"github.com/kaz/rusuden/recognizer/gcp"
)

var (
	GCS_BUCKET     = os.Getenv("GCS_BUCKET")
	MG_SIGNING_KEY = os.Getenv("MG_SIGNING_KEY")
	MG_API_KEY     = os.Getenv("MG_API_KEY")
	MG_SENDER      = os.Getenv("MG_SENDER")
	MG_RECIPIENT   = os.Getenv("MG_RECIPIENT")

	recognizer = gcp.New(GCS_BUCKET)
	arrival    = mgArr.New(MG_SIGNING_KEY)
	departure  = mgDep.New(MG_API_KEY, MG_SENDER, MG_RECIPIENT)
)

func Handle(w http.ResponseWriter, r *http.Request) {
	if err := handle(r); err != nil {
		log.Printf("handling error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handle(r *http.Request) error {
	data, err := arrival.Arrive(r)
	if err != nil {
		return fmt.Errorf("arrival.Arrive failed: %w", err)
	}

	text, err := recognizer.Recognize(r.Context(), data)
	if err != nil {
		return fmt.Errorf("recognizer.Recognize failed: %w", err)
	}

	err = departure.Depart(r.Context(), text)
	if err != nil {
		return fmt.Errorf("departure.Depart failed: %w", err)
	}

	return nil
}
