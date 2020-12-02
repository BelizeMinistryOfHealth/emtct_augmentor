package labs

import (
	"fmt"
	"time"
)

func (d *Labs) FindInfantSyphilisScreenings(infantId int, birthDate time.Time) ([]SyphilisScreening, error) {
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
	WHERE 
		t.test_id=1 AND p.patient_id=$1
		AND tr.order_received_by_lab_time < ($2::date + '2 year'::interval);
`
	var testRequests []testRequestItem
	dob := birthDate.Format(layoutISO)
	rows, err := d.AcsisDb.Query(stmt, infantId, dob)
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
	var labResults []LabResult
	//for t := range testRequests {
	//	testResults, err := d.findTestResults(testRequests[t])
	//	if err != nil {
	//		return nil, fmt.Errorf("")
	//	}
	//	for _, i := range testResults {
	//		result := models.LabResult{
	//			Id:                     i.Id,
	//			PatientId:              infantId,
	//			TestResult:             i.TestResult,
	//			TestName:               fmt.Sprintf("%s - %s", i.TestName, i.TestLabel),
	//			TestRequestItemId:      i.TestRequestItemId,
	//			DateSampleTaken:        nil,
	//			ResultDate:             nil,
	//			ReleasedTime:           testRequests[t].ReleasedTime,
	//			DateOrderReceivedByLab: testRequests[t].DateOrderReceivedByLab,
	//		}
	//		labResults = append(labResults, result)
	//	}
	//}

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
	var screenings []SyphilisScreening
	for _, l := range labResults {
		s := SyphilisScreening{
			Id:                 l.Id,
			PatientId:          l.PatientId,
			TestName:           l.TestName,
			ScreeningDate:      *l.DateOrderReceivedByLab,
			DateResultReceived: l.ReleasedTime,
			DateSampleTaken:    l.DateSampleTaken,
			Result:             l.TestResult,
			Timely:             NotAvailable,
		}
		screenings = append(screenings, s)
	}
	return screenings, nil
}
