package db

import (
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"moh.gov.bz/mch/emtct/internal/models"
)

type testRequestItem struct {
	PatientId              int
	EncounterId            int
	TestRequestId          int
	TestRequestItemId      int
	TestName               string
	TestResult             string
	ReleasedTime           *time.Time
	DateOrderReceivedByLab *time.Time
}

// findCurrentTestRequestItems finds all test requests in a given encounter. This is used when
// searching for a pregnant woman's test results during pregnancy.
func (d *AcsisDb) findCurrentTestRequestItems(patientId, encounterId int, ancBeginDate *time.Time) ([]testRequestItem, error) {
	stmt := `SELECT p.patient_id,
                    e.encounter_id,
                    tri.test_request_item_id,
                    tr.test_request_id,
       				tri.released_time,
       				tr.order_received_by_lab_time,
                    t.name
             FROM acsis_hc_patients p
             INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id 
                                                      --AND encounter_type IN ('M', 'B')
             INNER JOIN acsis_lab_test_requests tr ON tr.encounter_id=e.encounter_id
             INNER JOIN acsis_lab_test_request_items tri ON tr.test_request_id=tri.test_request_id
             INNER JOIN acsis_lab_tests t ON tri.test_id=t.test_id
             WHERE p.patient_id=$1 AND e.encounter_id=$2 AND $3 < tr.order_received_by_lab_time`
	var testRequests []testRequestItem
	rows, err := d.Query(stmt, patientId, encounterId, ancBeginDate.Format(layoutISO))
	if err != nil {
		return nil, fmt.Errorf("error retrieving test request items from acsis: %+v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t testRequestItem
		err := rows.Scan(&t.PatientId,
			&t.EncounterId,
			&t.TestRequestItemId,
			&t.TestRequestId,
			&t.ReleasedTime,
			&t.DateOrderReceivedByLab,
			&t.TestName)
		if err != nil {
			return nil, fmt.Errorf("error scanning test request item from acsis: %+v", err)
		}
		testRequests = append(testRequests, t)
	}

	return testRequests, nil
}

type testResult struct {
	Id                int
	PatientId         int
	TestRequestId     int
	TestRequestItemId int
	EncounterId       int
	TestName          string
	TestResult        string
	TestLabel         string
}

func (d *AcsisDb) findTestResults(ri testRequestItem) ([]testResult, error) {
	stmt := `
	SELECT 
	    a.test_result_id,
		altr.test_request_id,
		altri.test_request_item_id,
		e.encounter_id,
		alt.name as test,
		aludli.name as result,
		a.label
	FROM acsis_hc_patients p
		INNER JOIN acsis_adt_encounters AS e ON e.patient_id = p.patient_id
		INNER JOIN acsis_lab_test_requests altr on e.encounter_id = altr.encounter_id 
-- 		                                               AND encounter_type='M'
		INNER JOIN acsis_lab_test_request_items altri on altr.test_request_id = altri.test_request_id
		INNER JOIN acsis_lab_tests alt on altri.test_id = alt.test_id
		INNER JOIN acsis_lab_test_request_results_collected altrrc on altri.test_request_item_id = altrrc.test_request_item_id
		INNER JOIN acsis_lab_test_results a on altrrc.test_result_id = a.test_result_id
		INNER JOIN acsis_lab_user_defined_list_items aludli on altrrc.user_defined_list_value = aludli.user_defined_list_item_id
	WHERE p.patient_id=$1
		AND altr.test_request_id=$2
		AND e.active IS TRUE
	ORDER BY altr.last_modified_time DESC;
`
	var results []testResult
	rows, err := d.Query(stmt, ri.PatientId, ri.TestRequestId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving test results for a test request items from acsis: %+v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var r testResult
		err := rows.Scan(
			&r.Id,
			&r.TestRequestId,
			&r.TestRequestItemId,
			&r.EncounterId,
			&r.TestName,
			&r.TestResult,
			&r.TestLabel)
		if err != nil {
			return nil, fmt.Errorf("error scanning test results when fetching test results from acsis: %+v", err)
		}
		results = append(results, r)
	}
	return results, nil
}

type testSample struct {
	TestSampleId      sql.NullInt32
	CollectedTime     *time.Time
	TestRequestItemId int
	TestRequestId     int
}

func (d *AcsisDb) findTestSamples(tr testRequestItem) (*testSample, error) {
	stmt := `
	SELECT  
		alts.test_sample_id,
		alts.collected_time
	FROM acsis_lab_test_request_specimen_types altrst
		INNER JOIN acsis_lab_test_request_items altri ON altri.test_request_item_id=$1
		LEFT JOIN acsis_lab_test_samples alts ON alts.test_request_specimen_type_id=altrst.test_request_specimen_type_id
	WHERE altrst.test_request_id=$2
	ORDER BY alts.collected_time
	LIMIT 1;
`
	row := d.QueryRow(stmt, tr.TestRequestItemId, tr.TestRequestId)
	var s testSample
	err := row.Scan(&s.TestSampleId, &s.CollectedTime)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		s.TestRequestId = tr.TestRequestId
		s.TestRequestItemId = tr.TestRequestItemId
		return &s, nil
	default:
		return nil, fmt.Errorf("error retrieving lab sample info from acsis: %+v", err)
	}
}

func assignSamplesToResults(results []models.LabResult, samples []testSample) []models.LabResult {
	for _, s := range samples {
		for _, r := range results {
			if s.TestRequestItemId == r.TestRequestItemId {
				r.DateSampleTaken = s.CollectedTime
			}
		}
	}
	// Deduplicate results
	var r []models.LabResult
	for _, result := range results {
		index := findLabResultIndex(r, result.Id)
		if index < 0 {
			r = append(r, result)
		}
	}
	return r
}

func findLabResultIndex(rs []models.LabResult, id int) int {
	for i, v := range rs {
		if v.Id == id {
			return i
		}
	}
	return -1
}

// FindTestsDuringPregnancy returns all the tests conducted during a woman's pregnancy.
// Since we also need the date the sample was collected, the query gets more complicated.
// So we have to issue multiple queries to retrieve separate parts of the information.
// 0. Find the latest anc encounter.
// 1. Find all test request items.
// 2. Find test results for each test request item
// 3. Find the samples for each test request item
// 4. Create the response that will merge the data from all these queries.
func (d *AcsisDb) FindLabTestsDuringPregnancy(patientId int) ([]models.LabResult, error) {
	v, err := d.FindObstetricDetails(patientId)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving lab tests during pregnancy details from acsis: %+v", err)
	}
	anc, err := d.FindLatestAntenatalEncounter(patientId, v.Lmp)
	if err != nil {
		return nil, fmt.Errorf("rerror retrieving antenatal encounter when retrieving lab tests during pregnancy: %+v", err)
	}
	encounterId := anc.Id
	testItems, err := d.findCurrentTestRequestItems(patientId, encounterId, anc.BeginDate)
	if err != nil {
		return nil, fmt.Errorf("error finding current test request items from acsis when retrieving lab tests during pregnancy: %+v", err)
	}
	log.WithFields(log.Fields{"testItems": testItems}).Info("test reuqest Items")
	var labResults []models.LabResult
	for _, t := range testItems {
		testResults, err := d.findTestResults(t)
		if err != nil {
			return nil, fmt.Errorf("error finding test results for a test request item when retrieving lab tests during pregnancy from acsis: %+v", err)
		}
		for _, r := range testResults {
			result := models.LabResult{
				Id:                     r.Id,
				PatientId:              patientId,
				TestName:               fmt.Sprintf("%s - %s", r.TestName, r.TestLabel),
				TestResult:             r.TestResult,
				TestRequestItemId:      r.TestRequestItemId,
				ReleasedTime:           t.ReleasedTime,
				DateOrderReceivedByLab: t.DateOrderReceivedByLab,
				DateSampleTaken:        nil,
				ResultDate:             nil,
			}
			labResults = append(labResults, result)
		}
	}
	var testSamples []testSample
	for _, t := range testItems {
		sample, err := d.findTestSamples(t)
		if err != nil {
			return nil, fmt.Errorf("error finding test samples from when retrieving lab tests during prengnacy from acsis: %+v", err)
		}
		if sample != nil {
			testSamples = append(testSamples, *sample)
		}

	}
	results := assignSamplesToResults(labResults, testSamples)
	return results, nil
}

func (d *AcsisDb) FindInfantSyphilisScreenings(infantId int, birthDate time.Time) ([]models.SyphilisScreening, error) {
	stmt := `
SELECT 
	p.patient_id,
	e.encounter_id,
	tri.test_request_item_id,
	tr.test_request_id,
	tri.released_time,
	tr.order_received_by_lab_time,
	t.name
FROM acsis_hc_patients p
 	INNER JOIN acsis_adt_encounters e ON p.patient_id=e.patient_id AND encounter_type='M'
 	INNER JOIN acsis_lab_test_requests tr ON tr.encounter_id=e.encounter_id
 	INNER JOIN acsis_lab_test_request_items tri ON tr.test_request_id=tri.test_request_id
 	INNER JOIN acsis_lab_tests t ON tri.test_id=t.test_id
WHERE t.test_id=1 AND p.patient_id=$1
	AND tr.order_received_by_lab_time < ($2::date + '2 year'::interval);
`
	var testRequests []testRequestItem
	dob := birthDate.Format(layoutISO)
	rows, err := d.Query(stmt, infantId, dob)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("")
	}
	for rows.Next() {
		var t testRequestItem
		err := rows.Scan(&t.PatientId, &t.EncounterId, &t.TestRequestItemId, &t.TestRequestId, &t.ReleasedTime, &t.DateOrderReceivedByLab, &t.TestName)
		if err != nil {
			return nil, fmt.Errorf("")
		}
		testRequests = append(testRequests, t)
	}
	var labResults []models.LabResult
	for t := range testRequests {
		testResults, err := d.findTestResults(testRequests[t])
		if err != nil {
			return nil, fmt.Errorf("")
		}
		for _, i := range testResults {
			result := models.LabResult{
				Id:                     i.Id,
				PatientId:              infantId,
				TestResult:             i.TestResult,
				TestName:               fmt.Sprintf("%s - %s", i.TestName, i.TestLabel),
				TestRequestItemId:      i.TestRequestItemId,
				DateSampleTaken:        nil,
				ResultDate:             nil,
				ReleasedTime:           testRequests[t].ReleasedTime,
				DateOrderReceivedByLab: testRequests[t].DateOrderReceivedByLab,
			}
			labResults = append(labResults, result)
		}
	}

	var testSamples []testSample
	for _, t := range testRequests {
		sample, err := d.findTestSamples(t)
		if err != nil {
			return nil, fmt.Errorf("error finding test samples from when retrieving lab tests during prengnacy from acsis: %+v", err)
		}
		if sample != nil {
			testSamples = append(testSamples, *sample)
		}

	}

	labResults = assignSamplesToResults(labResults, testSamples)
	var screenings []models.SyphilisScreening
	for _, l := range labResults {
		s := models.SyphilisScreening{
			Id:                 l.Id,
			PatientId:          l.PatientId,
			TestName:           l.TestName,
			ScreeningDate:      *l.DateOrderReceivedByLab,
			DateResultReceived: l.ReleasedTime,
			DateSampleTaken:    l.DateSampleTaken,
			Result:             l.TestResult,
			Timely:             models.NotAvailable,
		}
		screenings = append(screenings, s)
	}
	return screenings, nil
}
