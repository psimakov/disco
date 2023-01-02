package disco

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/mchmarny/disco/pkg/disco"
	"github.com/mchmarny/disco/pkg/object"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
)

// DiscoHandler is the HTTP handler for disco service.
func (h *Handler) DiscoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	log.Println("preparing discovery...")

	requestContextID := uuid.NewString()
	dir, err := makeFolder(requestContextID)
	if err != nil {
		writeError(w, errors.Wrapf(err, "error creating context from: %s", requestContextID))
		return
	}
	defer func() {
		if err = os.RemoveAll(dir); err != nil {
			log.Printf("error deleting context: %s\n", dir)
		}
	}()

	reportName := fmt.Sprintf("%s.json", time.Now().UTC().Format("2006-01-02T15-04-05"))
	reportPath := path.Join(dir, reportName)
	query := &types.VulnsQuery{
		SimpleQuery: types.SimpleQuery{
			OutputPath: reportPath,
			OutputFmt:  types.JSONFormat,
		},
	}

	if err := disco.DiscoverVulns(r.Context(), query); err != nil {
		writeError(w, errors.Wrap(err, "error validating"))
		return
	}

	if err := object.Save(r.Context(), h.bucket, reportName, reportPath); err != nil {
		writeError(w, errors.Wrapf(err, "error writing content to: %s/%s",
			h.bucket, reportName))
		return
	}

	if err := disco.MeterVulns(r.Context(), h.counter, reportPath); err != nil {
		writeError(w, errors.Wrap(err, "error metering vulnerabilities"))
		return
	}

	writeMessage(w, fmt.Sprintf("gs://%s/%s", h.bucket, reportName))
}