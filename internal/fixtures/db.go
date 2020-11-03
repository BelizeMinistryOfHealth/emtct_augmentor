package fixtures

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"moh.gov.bz/mch/emtct/internal/db"
	"moh.gov.bz/mch/emtct/internal/models"
)

var PatientIds = []string{"1111120", "1111121"}

const (
	layoutISO = "2006-01-02"
)

func ClearTable(table string, db db.EmtctDb) error {
	stmt := fmt.Sprintf("DELETE FROM %s", table)
	_, err := db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error deleting table(%s) content: %+v", table, err)
	}
	return nil
}

func SamplePatients(db db.EmtctDb) error {

	for i, id := range PatientIds {
		patient := models.Patient{
			Id:               id,
			FirstName:        fmt.Sprintf("First Name - %d", i),
			MiddleName:       "",
			LastName:         fmt.Sprintf("Last Name - %d", i),
			Dob:              time.Date(1992, time.October, 1, 0, 0, 0, 0, time.UTC),
			Ssn:              uuid.New().String(),
			CountryOfBirth:   "Belize",
			DistrictAddress:  "Cayo",
			CommunityAddress: "Las Flores",
			Education:        "High School",
			Ethnicity:        "Mestizo",
			Hiv:              false,
			NextOfKin:        "Jon Doe",
			NextOfKinPhone:   "6539333",
		}
		err := db.CreatePatient(patient)
		if err != nil {
			return fmt.Errorf("failure creating sample patient: %+v", err)
		}
	}
	return nil
}

func SampleDiagnoses(db db.EmtctDb) error {
	patientId := PatientIds[0]
	ds := []string{"2009-02-03", "2010-11-15", "2011-04-12", "2020-10-23"}
	var dates []time.Time
	for _, d := range ds {
		date, _ := time.Parse(layoutISO, d)
		dates = append(dates, date)
	}

	diagnoses := []models.Diagnosis{
		{
			Id:        "1",
			PatientId: patientId,
			Date:      dates[0],
			Name:      "common cold",
		},
		{
			Id:        "2",
			PatientId: patientId,
			Date:      dates[1],
			Name:      "seasonal flu",
		},
		{
			Id:        "3",
			PatientId: patientId,
			Date:      dates[2],
			Name:      "rash",
		},
		{
			Id:        "4",
			PatientId: patientId,
			Date:      dates[3],
			Name:      "common cold",
		},
	}
	for _, d := range diagnoses {
		err := db.CreateDiagnosis(d)
		if err != nil {
			return fmt.Errorf("error creating sample diagnosis: %+v", d)
		}
	}
	return nil
}

func SampleObstetricHistory(db db.EmtctDb) error {
	history := []models.ObstetricHistory{
		{
			Id:             "1",
			PatientId:      PatientIds[0],
			Date:           time.Date(2012, time.May, 1, 0, 0, 0, 0, time.UTC),
			ObstetricEvent: "Miscarriage",
		},
		{
			Id:             "2",
			PatientId:      PatientIds[0],
			Date:           time.Date(2015, time.August, 27, 0, 0, 0, 0, time.UTC),
			ObstetricEvent: "Live Born",
		},
	}
	for _, h := range history {
		err := db.CreateObstetricHistory(h)
		if err != nil {
			return fmt.Errorf("error creating sample obstetric history: %+v", err)
		}
	}
	return nil
}

func SamplePregnancies(db db.EmtctDb) error {
	patientId, _ := strconv.Atoi(PatientIds[0])
	pregnancy := models.PregnancyVitals{
		Id:                   1,
		PatientId:            patientId,
		GestationalAge:       4,
		Para:                 10,
		Cs:                   false,
		PregnancyOutcome:     "",
		DiagnosisDate:        time.Date(2020, time.August, 3, 0, 0, 0, 0, time.UTC),
		Planned:              false,
		AgeAtLmp:             28,
		Lmp:                  time.Date(2020, time.July, 6, 0, 0, 0, 0, time.UTC),
		Edd:                  time.Date(2021, time.April, 4, 0, 0, 0, 0, time.UTC),
		DateOfBooking:        time.Date(2020, time.August, 3, 0, 0, 0, 0, time.UTC),
		PrenatalCareProvider: "Public",
		TotalChecks:          2,
	}
	err := db.CreatePregnancy(pregnancy)
	if err != nil {
		return fmt.Errorf("error inserting pregnancy: %+v", err)
	}
	return nil
}

func SampleLabResults(db db.EmtctDb) error {
	patientId, _ := strconv.Atoi(PatientIds[0])
	labResults := []models.LabResult{
		{
			Id:              1,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "Hb",
			DateSampleTaken: time.Date(2020, time.September, 10, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.September, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              2,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "Urinalysis",
			DateSampleTaken: time.Date(2020, time.September, 10, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.September, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              3,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "Hepatitis B",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              4,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "HIV",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              5,
			PatientId:       patientId,
			TestResult:      "120",
			TestName:        "CD4 Count",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              6,
			PatientId:       patientId,
			TestResult:      "0",
			TestName:        "Viral Load",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:              7,
			PatientId:       patientId,
			TestResult:      "Negative",
			TestName:        "Syphilis",
			DateSampleTaken: time.Date(2020, time.June, 30, 0, 0, 0, 0, time.UTC),
			ResultDate:      time.Date(2020, time.July, 3, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, l := range labResults {
		err := db.CreateLabResult(l)
		if err != nil {
			return fmt.Errorf("error inserting lab result: %+v", err)
		}
	}
	return nil
}

func SampleHomeVisits(db db.EmtctDb) error {
	patientId, _ := strconv.Atoi(PatientIds[0])
	homeVisits := []models.HomeVisit{
		{
			Id:          uuid.New().String(),
			PatientId:   patientId,
			Reason:      "Random",
			Comments:    "Patient's vitals are normal",
			CreatedAt:   time.Date(2020, time.September, 29, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   nil,
			CreatedBy:   "nurse@health.gov.bz",
			UpdatedBy:   nil,
			DateOfVisit: time.Date(2020, time.September, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			Id:          uuid.New().String(),
			PatientId:   patientId,
			Reason:      "Periodic",
			Comments:    "All vitals were normal. Patient was given information on breast feeding.",
			CreatedAt:   time.Date(2020, time.October, 16, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   nil,
			CreatedBy:   "nurse@health.gov.bz",
			UpdatedBy:   nil,
			DateOfVisit: time.Date(2020, time.October, 16, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, h := range homeVisits {
		err := db.CreateHomeVisit(h)
		if err != nil {
			return fmt.Errorf("error inserting sample home visit: %+v", err)
		}
	}
	return nil
}
