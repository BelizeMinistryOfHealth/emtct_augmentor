package models

import (
	"time"
)

type Patient struct {
	Id               string    `json:"patientId"`
	PregnancyId      string    `json:"pregnancyId"`
	FirstName        string    `json:"firstName"`
	MiddleName       string    `json:"middleName"`
	LastName         string    `json:"lastName"`
	Dob              time.Time `json:"dob"`
	Ssn              string    `json:"ssn"`
	CountryOfBirth   string    `json:"countryOfBirth"`
	DistrictAddress  string    `json:"district"`
	CommunityAddress string    `json:"community"`
	Address          string    `json:"address"`
	Education        string    `json:"education"`
	Ethnicity        string    `json:"ethnicity"`
	Hiv              bool      `json:"hiv"`
	HivDiagnosisDate time.Time `json:"hivDiagnosisDate"`
	NextOfKin        string    `json:"nextOfKin"`
	NextOfKinPhone   string    `json:"nextOfKinPhone"`
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

func FindCurrentPregnancy(ps []PregnancyVitals) *PregnancyVitals {
	for _, v := range ps {
		if v.Edd.After(time.Now()) {
			return &v
		}
	}
	return nil
}

type LabResult struct {
	Id                     int        `json:"id"`
	PatientId              int        `json:"patientId"`
	TestResult             string     `json:"testResult"`
	TestName               string     `json:"testName"`
	TestRequestItemId      int        `json:"testRequestItemId"`
	DateSampleTaken        *time.Time `json:"dateSampleTaken"`
	ResultDate             *time.Time `json:"resultDate"`
	ReleasedTime           *time.Time `json:"releasedTime"`
	DateOrderReceivedByLab *time.Time `json:"dateOrderReceivedByLab"`
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

type HivScreening struct {
	Id                     string     `json:"id"`
	PatientId              int        `json:"patientId"`
	MchEncounterId         int        `json:"mchEncounterId"`
	TestName               string     `json:"testName"`
	ScreeningDate          time.Time  `json:"screeningDate"`
	DateSampleReceivedAtHq *time.Time `json:"dateSampleReceivedAtHq,omitEmpty"`
	SampleCode             string     `json:"sampleCode"`
	DateSampleShipped      *time.Time `json:"dateSampleShipped"`
	Destination            string     `json:"destination"`
	DateResultReceived     *time.Time `json:"dateResultReceived,omitEmpty"`
	DateSampleTaken        *time.Time `json:"dateSampleTaken,omitEmpty"`
	DueDate                *time.Time `json:"dueDate,omitEmpty"`
	Result                 string     `json:"result"`
	DateResultShared       *time.Time `json:"dateResultShared,omitEmpty"`
	Timely                 bool       `json:"timely"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              *time.Time `json:"updatedAt"`
	CreatedBy              string     `json:"createdBy"`
	UpdatedBy              *string    `json:"updatedBy"`
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

// IsHivScreeningTimely indicates if an hiv screening was done in a timely manner.
// The timeliness depends on the type of test and when the sample was taken:
// PCR 1: sample must be taken 3 days or less after birth.
// PCR 2: sample must be taken no later than 6 weeks after birth
// PCR 3: sample must be taken no later than 90 days after birth
// ELISA: sample must be taken no longer than 18 months after birth
func IsHivScreeningTimely(birth Birth, testName string, dateSampleTaken time.Time) bool {
	diff := dateSampleTaken.Sub(birth.BirthDate).Hours() / 24
	switch testName {
	case "PCR 1":
		return diff < 4
	case "PCR 2":
		return diff < (6 * 7)
	case "PCR 3":
		return diff < 91
	case "ELISA":
		return diff <= (18 * 7 * 4)
	default:
		return false
	}
}

// HivScreeningDueDate calculates the due date for taking a sample for an HIV screening.
// PCR 1: sample must be taken 3 days or less after birth.
// PCR 2: sample must be taken no later than 6 weeks after birth
// PCR 3: sample must be taken no later than 90 days after birth
// ELISA: sample must be taken no longer than 18 months after birth
func HivScreeningDueDate(testName string, birthDate time.Time) time.Time {
	switch testName {
	case "PCR 1":
		return birthDate.AddDate(0, 0, 3)
	case "PCR 2":
		return birthDate.AddDate(0, 0, 42)
	case "PCR 3":
		return birthDate.AddDate(0, 0, 90)
	default:
		return birthDate.AddDate(0, 18, 0)
	}
}
