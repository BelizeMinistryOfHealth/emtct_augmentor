package models

import (
	"time"
)

type Patient struct {
	Id               string     `json:"patientId"`
	PregnancyId      string     `json:"pregnancyId"`
	FirstName        string     `json:"firstName"`
	MiddleName       string     `json:"middleName"`
	LastName         string     `json:"lastName"`
	Dob              time.Time  `json:"dob"`
	Ssn              string     `json:"ssn"`
	CountryOfBirth   string     `json:"countryOfBirth"`
	DistrictAddress  string     `json:"district"`
	CommunityAddress string     `json:"community"`
	Address          string     `json:"address"`
	Education        string     `json:"education"`
	Ethnicity        string     `json:"ethnicity"`
	Hiv              bool       `json:"hiv"`
	HivDiagnosisDate *time.Time `json:"hivDiagnosisDate"`
	NextOfKin        string     `json:"nextOfKin"`
	NextOfKinPhone   string     `json:"nextOfKinPhone"`
}

type PatientBasicInfo struct {
	Id         string    `json:"patientId"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	Dob        time.Time `json:"dob"`
	Ssn        string    `json:"ssn"`
}

type ObstetricHistory struct {
	Id             string    `json:"id"`
	PatientId      string    `json:"patientId"`
	Date           time.Time `json:"date"`
	ObstetricEvent string    `json:"event"`
}

type Diagnosis struct {
	Id        string    `json:"id"`
	PatientId string    `json:"patientId"`
	Date      time.Time `json:"date"`
	Name      string    `json:"name"`
}

type PregnancyVitals struct {
	Id                   int        `json:"id"`
	PatientId            int        `json:"patientId"`
	GestationalAge       int        `json:"gestationalAge"`
	Para                 int        `json:"para"`
	Cs                   int        `json:"cs"`
	PregnancyOutcome     string     `json:"pregnancyOutcome"`
	DiagnosisDate        *time.Time `json:"diagnosisDate"`
	Planned              bool       `json:"planned"`
	AgeAtLmp             int        `json:"ageAtLmp"`
	Lmp                  *time.Time `json:"lmp"`
	Edd                  time.Time  `json:"edd"`
	DateOfBooking        *time.Time `json:"dateOfBooking"`
	PrenatalCareProvider string     `json:"prenatalCareProvider"`
	TotalChecks          int        `json:"totalChecks"`
	ApgarFirstMinute     int        `json:"apgarFirstMinute"`
	ApgarFifthMinute     int        `json:"apgarFifthMinute"`
	BirthStatus          string     `json:"birthStatus"`
}

type LabResult struct {
	Id                     int        `json:"id"`
	PatientId              int        `json:"patientId"`
	TestResult             string     `json:"testResult"`
	TestName               string     `json:"testName"`
	TestRequestId          int        `json:"testRequestId"`
	TestRequestItemId      int        `json:"testRequestItemId"`
	DateSampleTaken        *time.Time `json:"dateSampleTaken"`
	ResultDate             *time.Time `json:"resultDate"`
	ReleasedTime           *time.Time `json:"releasedTime"`
	DateOrderReceivedByLab *time.Time `json:"dateOrderReceivedByLab"`
}

type LabTest struct {
	Id                     string     `json:"id"`
	PatientId              int        `json:"patientId"`
	TestResult             string     `json:"testResult"`
	TestName               string     `json:"testName"`
	TestRequestId          int        `json:"testRequestId"`
	TestRequestItemId      int        `json:"testRequestItemId"`
	DateSampleTaken        *time.Time `json:"dateSampleTaken"`
	ResultDate             *time.Time `json:"resultDate"`
	ReleasedTime           *time.Time `json:"releasedTime"`
	DateOrderReceivedByLab *time.Time `json:"dateOrderReceivedByLab"`
}

type Pregnancy struct {
	PatientId   int
	PregnancyId int
	Lmp         *time.Time
	Edd         *time.Time
	EndTime     *time.Time
}

func (p *Pregnancy) Index(vs []Pregnancy) int {
	for i, v := range vs {
		if v.PregnancyId == p.PregnancyId {
			return i
		}
	}
	return -1
}

func (p *Pregnancy) Include(vs []Pregnancy) bool {
	return p.Index(vs) >= 0
}

// FindLabResultsBetweenDates finds lab results between a range of two dates.
// It is specifically for finding lab results during a pregnancy period.
func FindLabResultsBetweenDates(labResults []LabResult, lmp time.Time) []LabResult {
	var results []LabResult
	for _, l := range labResults {
		if l.DateSampleTaken.After(lmp) && l.DateSampleTaken.Before(lmp.Add(time.Hour*24*30*9)) {
			results = append(results, l)
		}
	}

	return results
}

type HomeVisit struct {
	Id             string     `json:"id"`
	PatientId      int        `json:"patientId"`
	MchEncounterId int        `json:"mchEncounterId"`
	Reason         string     `json:"reason"`
	Comments       string     `json:"comments"`
	DateOfVisit    time.Time  `json:"dateOfVisit"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	CreatedBy      string     `json:"createdBy"`
	UpdatedBy      *string    `json:"updatedBy"`
}

type SampleTimeliness string

const (
	Timely       SampleTimeliness = "Timely"
	NotTimely    SampleTimeliness = "NotTimely"
	NotAvailable SampleTimeliness = "N/A"
)

type SyphilisScreening struct {
	Id                 int              `json:"id"`
	PatientId          int              `json:"patientId"`
	TestName           string           `json:"testName"`
	ScreeningDate      time.Time        `json:"screeningDate"`
	DateResultReceived *time.Time       `json:"dateResultReceived,omitEmpty"`
	DateSampleTaken    *time.Time       `json:"dateSampleTaken,omitEmpty"`
	DueDate            *time.Time       `json:"dueDate,omitEmpty"`
	Result             string           `json:"result"`
	DateResultShared   *time.Time       `json:"dateResultShared,omitEmpty"`
	Timely             SampleTimeliness `json:"timely"`
}

type ContraceptiveUsed struct {
	Id             string     `json:"id"`
	PatientId      int        `json:"patientId"`
	MchEncounterId int        `json:"mchEncounterId"`
	Contraceptive  string     `json:"contraceptive"`
	Comments       string     `json:"comments"`
	DateUsed       time.Time  `json:"dateUsed"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	CreatedBy      string     `json:"createdBy"`
	UpdatedBy      *string    `json:"updatedBy"`
}

type HospitalAdmission struct {
	Id             string     `json:"id"`
	PatientId      int        `json:"patientId"`
	MchEncounterId int        `json:"mchEncounterId"`
	DateAdmitted   time.Time  `json:"dateAdmitted"`
	Facility       string     `json:"facility"`
	Reason         string     `json:"reason"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	CreatedBy      string     `json:"createdBy"`
	UpdatedBy      *string    `json:"updatedBy"`
}

type AntenatalEncounter struct {
	Id                    int        `json:"id"`
	PatientId             int        `json:"patientId"`
	MchEncounterDetailsId int        `json:"mchEncounterDetailsId"`
	EstimatedDeliveryDate *time.Time `json:"estimatedDeliveryDate"`
	BeginDate             *time.Time `json:"beginDate"`
	GestationalAge        int        `json:"gestationalAge"`
	NumberAntenatalVisits int        `json:"numberAntenatalVisits"`
}

type Prescription struct {
	Id             int       `json:"id"`
	PatientId      int       `json:"patientId"`
	TotalDoses     int       `json:"totalDoses"`
	Pharmaceutical string    `json:"pharmaceutical"`
	Frequency      string    `json:"frequency"`
	Strength       string    `json:"strength"`
	Comments       string    `json:"comments"`
	PrescribedTime time.Time `json:"prescribedTime"`
}

// SyphilisTreatment describes the treatment given to a patient's contact.
// It is very similar to the Prescription struct, but we do not capture the person's
// name.
type SyphilisTreatment struct {
	Id         string     `json:"id"`
	PatientId  int        `json:"patientId"`
	Medication string     `json:"medication"`
	Dosage     string     `json:"dosage"`
	Comments   string     `json:"comments"`
	Date       time.Time  `json:"date"`
	CreatedBy  string     `json:"createdBy"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedBy  string     `json:"updatedBy"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}

type InfantDiagnoses struct {
	DiagnosisId int       `json:"diagnosisId"`
	PatientId   int       `json:"patientId"`
	Diagnosis   string    `json:"diagnosis"`
	Doctor      string    `json:"doctor"`
	Comments    string    `json:"comments"`
	Date        time.Time `json:"date"`
}

type Birth struct {
	PatientId   int
	BirthStatus string
	Date        time.Time
	BirthDate   time.Time
}

type Person struct {
	PatientId  int        `json:"patientId"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	MiddleName string     `json:"middleName"`
	Dob        *time.Time `json:"dob"`
}

type Infant struct {
	Infant Person `json:"infant"`
	Mother Person `json:"mother"`
}

type ContactTracing struct {
	Id         string     `json:"id"`
	PatientId  int        `json:"patientId"`
	Test       string     `json:"test"`
	TestResult string     `json:"testResult"`
	Comments   string     `json:"comments"`
	Date       time.Time  `json:"date"`
	CreatedBy  string     `json:"createdBy"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedBy  string     `json:"updatedBy"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
