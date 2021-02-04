package reports

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"moh.gov.bz/mch/emtct/internal/business/data/infant"
	"moh.gov.bz/mch/emtct/internal/models"
	"time"
)

type MissingPcr struct {
	Screening infant.HivScreening `json:"screening"`
	Infant    models.Infant       `json:"infant"`
}

type hivScreeningByInfant struct {
	Screenings []infant.HivScreening
	InfantId   string
}

// index returns the index position of a string in an array of strings.
// Returns -1 if the string does not exist in the array.
func index(vs []hivScreeningByInfant, t string) int {
	for i, v := range vs {
		if v.InfantId == t {
			return i
		}
	}
	return -1
}

// include indicates if a string is present in an array.
func include(vs []hivScreeningByInfant, t string) bool {
	return index(vs, t) >= 0
}

// groupHivScreenings groups all HivScreenings by infant id
func groupHivScreenings(s []infant.HivScreening) []hivScreeningByInfant {
	var screenings []hivScreeningByInfant
	for _, x := range s {
		i := index(screenings, x.PatientId)
		if i < 0 {
			h := hivScreeningByInfant{
				Screenings: []infant.HivScreening{x},
				InfantId:   x.PatientId,
			}
			screenings = append(screenings, h)
		} else {
			sc := screenings[i]
			h := sc.Screenings
			h = append(h, x)
			sc.Screenings = h
			screenings[i] = sc
		}
	}
	return screenings
}

func (d *Reports) toMissingPcr(ctx context.Context, s []infant.HivScreening) ([]MissingPcr, error) {
	byInfantId := groupHivScreenings(s)
	ref := d.firestore.Client.Collection(d.collections.Infants)
	var missingPcrs []MissingPcr
	for _, i := range byInfantId {
		snap, err := ref.Doc(i.InfantId).Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("error fetching infant (%s) when: %w", i.InfantId, err)
		}
		var pt models.Infant
		if err := snap.DataTo(&pt); err != nil {
			return nil, fmt.Errorf("error transforming patient data for mother (%s): %w", i.InfantId, err)
		}
		for _, screening := range i.Screenings {
			missingPcrs = append(missingPcrs, MissingPcr{
				Screening: screening,
				Infant:    pt,
			})
		}
	}

	return missingPcrs, nil
}

func (d *Reports) MissingPcrs(ctx context.Context, year int) ([]MissingPcr, error) {
	// Find all the hiv screenings due on the particular year but have no dueDate
	date, err := time.Parse("2006-01-02", fmt.Sprintf("%d-01-01", year))
	if err != nil {
		return nil, fmt.Errorf("failed to create date for MssingPCrs report: %w", err)
	}
	ref := d.firestore.Client.Collection(d.collections.HivScreening)
	iter := ref.Where("dueDate", ">=", date).Where("screeningDate", "==", nil).Documents(ctx)
	var reports []infant.HivScreening
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error querying for missing pcrs: %w", err)
		}
		var r infant.HivScreening
		if err := doc.DataTo(&r); err != nil {
			return nil, fmt.Errorf("error transforming missing pcr report data: %w", err)
		}
		reports = append(reports, r)
	}
	missingPcrs, err := d.toMissingPcr(ctx, reports)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch missing pcrs: %w", err)
	}
	return missingPcrs, nil
}
