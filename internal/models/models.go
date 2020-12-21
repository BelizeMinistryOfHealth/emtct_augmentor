package models

import "time"

type Patient struct {
	ID               string           `json:"patientId" firestore:"id"`
	PregnancyId      string           `json:"pregnancyId" firestore:"pregnancyId"`
	FirstName        string           `json:"firstName" firestore:"firstName"`
	MiddleName       string           `json:"middleName" firestore:"middleName"`
	LastName         string           `json:"lastName" firestore:"lastName"`
	Dob              time.Time        `json:"dob" firestore:"dob"`
	Ssn              string           `json:"ssn" firestore:"ssn"`
	CountryOfBirth   string           `json:"countryOfBirth" firestore:"countryOfBirth"`
	DistrictAddress  string           `json:"district" firestore:"district"`
	CommunityAddress string           `json:"community" firestore:"community"`
	Address          string           `json:"address" firestore:"address"`
	Education        string           `json:"education" firestore:"education"`
	Ethnicity        string           `json:"ethnicity" firestore:"ethnicity"`
	Hiv              bool             `json:"hiv" firestore:"hiv"`
	HivDiagnosisDate *time.Time       `json:"hivDiagnosisDate" firestore:"hivDiagnosisDate"`
	NextOfKin        string           `json:"nextOfKin" firestore:"nextOfKin"`
	NextOfKinPhone   string           `json:"nextOfKinPhone" firestore:"nextOfKinPhone"`
	ObstetricEvents  []ObstetricEvent `json:"obstetricEvents,omitempty" firestore:"obstetricEvents"`
}

type ObstetricEvent struct {
	ID             string    `json:"id" firestore:"id"`
	PatientId      string    `json:"patientId" firestore:"patientId"`
	Date           time.Time `json:"date" firestore:"date"`
	ObstetricEvent string    `json:"event" firestore:"event"`
}

func ObstetricEventsByPatient(events []ObstetricEvent, id string) []ObstetricEvent {
	var evs []ObstetricEvent
	for _, e := range events {
		if e.PatientId == id {
			evs = append(evs, e)
		}
	}
	return evs
}

type Diagnosis struct {
	ID        int       `json:"id" firestore:"id"`
	PatientId string    `json:"patientId" firestore:"patientId"`
	Date      time.Time `json:"date" firestore:"date"`
	Name      string    `json:"name" firestore:"name"`
	Doctor    string    `json:"doctor" firestore:"doctor"`
	Comments  string    `json:"comments" firestore:"comments"`
}

type ObstetricDetails struct {
	Id               int        `json:"id" firestore:"id"`
	GestationalAge   int        `json:"gestationalAge" firestore:"gestationalAge"`
	Para             int        `json:"para" firestore:"para"`
	Cs               int        `json:"cs" firestore:"cs"`
	PregnancyOutcome string     `json:"pregnancyOutcome" firestore:"pregnancyOutcome"`
	DiagnosisDate    *time.Time `json:"diagnosisDate" firestore:"diagnosisDate"`
	Planned          bool       `json:"planned" firestore:"planned"`
	AgeAtLmp         int        `json:"ageAtLmp" firestore:"ageAtLmp"`
	Lmp              *time.Time `json:"lmp" firestore:"lmp"`
	Edd              time.Time  `json:"edd" firestore:"edd"`
	DateOfBooking    *time.Time `json:"dateOfBooking" firestore:"dateOfBooking"`
	ApgarFirstMinute int        `json:"apgarFirstMinute" firestore:"apgarFirstMinute"`
	ApgarFifthMinute int        `json:"apgarFifthMinute" firestore:"apgarFifthMinute"`
	BirthStatus      string     `json:"birthStatus" firestore:"birthStatus"`
}

type AntenatalEncounter struct {
	Id                    int        `json:"id" firestore:"id"`
	PatientId             string     `json:"patientId" firestore:"patientId"`
	MchEncounterDetailsId int        `json:"mchEncounterDetailsId" firestore:"mchEncounterDetailsId"`
	EstimatedDeliveryDate *time.Time `json:"estimatedDeliveryDate" firestore:"estimatedDeliveryDate"`
	DateOfBooking         *time.Time `json:"dateOfBooking" firestore:"dateOfBooking"`
	GestationalAge        int        `json:"gestationalAge" firestore:"gestationalAge"`
	NumberAntenatalVisits int        `json:"numberAntenatalVisits" firestore:"numberAntenatalVisits"`
	CareProvider          string     `json:"careProvider" firestore:"careProvider"`
	Facility              string     `json:"facility" firestore:"facility"`
	BirthStatus           string     `json:"birthStatus" firestore:"birthStatus"`
	ApgarFirstMinute      int        `json:"apgarFirstMinute" firestore:"apgarFirstMinute"`
	ApgarFifthMinute      int        `json:"apgarFifthMinute" firestore:"apgarFifthMinute"`
	DeliveryDate          *time.Time `json:"deliveryDate" firestore:"deliveryDate"`
}

type Pregnancy struct {
	PatientId        string             `json:"patientId" firestore:"patientId"`
	PregnancyId      int                `json:"id" firestore:"id"`
	Lmp              *time.Time         `json:"lmp" firestore:"lmp"`
	Edd              *time.Time         `json:"edd" firestore:"edd"`
	EndTime          *time.Time         `json:"endTime" firestore:"endTime"`
	ObstetricDetails ObstetricDetails   `json:"obstetricDetails" firestore:"obstetricDetails"`
	Anc              AntenatalEncounter `json:"anc" firestore:"anc"`
}

func PregnancyIndex(vs []Pregnancy, id int) int {
	for i, v := range vs {
		if v.PregnancyId == id {
			return i
		}
	}
	return -1
}

type LabResult struct {
	TestRequestItemId      int        `json:"testRequestItemId" firestore:"testRequestItemId"`
	LabResultId            int        `json:"labResultId" firestore:"labResultId"`
	PatientId              string     `json:"patientId" firestore:"patientId"`
	TestResult             string     `json:"testResult" firestore:"testResult"`
	TestName               string     `json:"testName" firestore:"testName"`
	TestRequestId          int        `json:"testRequestId" firestore:"testRequestId"`
	DateSampleTaken        *time.Time `json:"dateSampleTaken" firestore:"dateSampleTaken"`
	ResultDate             *time.Time `json:"resultDate" firestore:"resultDate"`
	ReleasedTime           *time.Time `json:"releasedTime" firestore:"releasedTime"`
	DateOrderReceivedByLab *time.Time `json:"dateOrderReceivedByLab" firestore:"dateOrderReceivedByLab"`
}

func LabResultIndex(vs []LabResult, testRequestItemId int) *LabResult {
	for _, v := range vs {
		if v.TestRequestItemId == testRequestItemId {
			return &v
		}
	}
	return nil
}

type Prescription struct {
	ID               int       `json:"id"`
	PatientId        string    `json:"patientId"`
	TotalDoses       int       `json:"totalDoses"`
	Pharmaceutical   string    `json:"pharmaceutical"`
	Frequency        string    `json:"frequency"`
	Strength         string    `json:"strength"`
	Comments         string    `json:"comments"`
	PrescribedTime   time.Time `json:"prescribedTime"`
	PharmaceuticalId int       `json:"pharmaceuticalId" firestore:"pharmaceuticalId"`
}

type Person struct {
	PatientId  int        `json:"patientId" firestore:"patientId"`
	FirstName  string     `json:"firstName" firestore:"firstName"`
	LastName   string     `json:"lastName" firestore:"lastName"`
	MiddleName string     `json:"middleName" firestore:"middleName"`
	Dob        *time.Time `json:"dob" firestore:"dob"`
}

type Infant struct {
	ID         string     `json:"id" firestore:"id"`
	FirstName  string     `json:"firstName" firestore:"firstName"`
	LastName   string     `json:"lastName" firestore:"lastName"`
	MiddleName string     `json:"middleName" firestore:"middleName"`
	Dob        *time.Time `json:"dob" firestore:"dob"`
	Mother     Person     `json:"mother" firestore:"mother"`
}
