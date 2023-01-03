package bq

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/mchmarny/disco/pkg/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func ImportLicenses(ctx context.Context, req *types.ImportRequest) error {
	if req == nil || req.FilePath == "" {
		return errors.Errorf("configured import request is required: %v", req)
	}
	var report types.ItemReport[types.LicenseReport]
	if err := types.UnmarshalFromFile(req.FilePath, &report); err != nil {
		return errors.Wrapf(err, "failed to unmarshal report from %s", req.FilePath)
	}

	rows := MakeLicenseRows(report.Items)
	if err := insert(ctx, req, rows); err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	log.Debug().Msgf("inserted %d rows into %s.%s.%s", len(rows), req.ProjectID, req.DatasetID, req.TableID)

	return nil
}

func MakeLicenseRows(in []*types.LicenseReport) []*LicenseRow {
	list := make([]*LicenseRow, 0)
	updated := time.Now().UTC().Format(time.RFC3339)

	for _, r := range in {
		for _, l := range r.Licenses {
			list = append(list, &LicenseRow{
				Image:   types.ParseImageNameFromDigest(r.Image),
				Sha:     types.ParseImageShaFromDigest(r.Image),
				Name:    l.Name,
				Source:  l.Source,
				Updated: updated,
			})
		}
	}

	return list
}

type LicenseRow struct {
	Image   string
	Sha     string
	Name    string
	Source  string
	Updated string
}

func (i *LicenseRow) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"image":   i.Image,
		"sha":     i.Sha,
		"name":    i.Name,
		"source":  i.Source,
		"updated": i.Updated,
	}, "", nil
}