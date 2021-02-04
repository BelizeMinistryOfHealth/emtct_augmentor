import _ from 'lodash';
/**
 * compilePcrs transforms a list of raw pcrs by keeping
 * all infant related pcrs in a test.
 * The API returns a row per PCR test, but we want to include all
 * PCRs in one row.
 * @param pcrs
 */
export const compilePcrs = (pcrs, infantId) => {
  //Key by infant's id
  const byInfantId = _.groupBy(pcrs, (p) => p.infant.id);
  //Now we can retrieve all the pcrs for an infant.
  const screenings = byInfantId[infantId];
  if (!screenings) {
    return {};
  }
  // The infant & mother info will be the same for all the records for a particular infant.
  // So we only need to get one screening to extract the infant & mother details.
  const baseData = screenings[0];
  const mergedScreenings = {
    id: baseData.infant.id,
    patientId: baseData.infant.id,
    infantName: `${baseData.infant.firstName} ${baseData.infant.lastName}`,
    infantDob: `${baseData.infant.dob}`,
    motherName: `${baseData.infant.mother.firstName} ${baseData.infant.mother.lastName}`,
    motherDob: `${baseData.infant.mother.dob}`,
    PCR1DueDate: screenings.find((s) => s.screening.testName === 'PCR 1')
      ?.screening?.dueDate,
    PCR1DateSampleTaken: screenings.find(
      (s) => s.screening.testName === 'PCR 1'
    )?.screening?.dateSampleTaken,
    PCR2DueDate: screenings.find((s) => s.screening.testName === 'PCR 2')
      ?.screening?.dueDate,
    PCR2DateSampleTaken: screenings.find(
      (s) => s.screening.testName === 'PCR 2'
    )?.screening?.dateSampleTaken,
    PCR3DueDate: screenings.find((s) => s.screening.testName === 'PCR 3')
      ?.screening?.dueDate,
    PCR3DateSampleTaken: screenings.find(
      (s) => s.screening.testName === 'PCR 3'
    )?.screening?.dateSampleTaken,
    ELISADueDate: screenings.find((s) => s.screening.testName === 'ELISA')
      ?.screening?.dueDate,
    ELISADateSampleTaken: screenings.find(
      (s) => s.screening.testName === 'ELISA'
    )?.screening?.dateSampleTaken,
  };

  return mergedScreenings;
};

/**
 * mergeScreenings take s a list of pcr screenings and merges all the test details
 * for infants into one row.
 * @param pcrs
 * @returns {({}|{motherDob: string, ELISADueDate: *, PCR1DueDate: *, patientId: *, ELISADateSampleTaken: *, motherName: string, infantName: string, infantDob: string, PCR1DateSampleTaken: *, PCR3DateSampleTaken: *, PCR3DueDate: *, PCR2DueDate: *, id: *, PCR2DateSampleTaken: *})[]}
 */
export const mergeScreenings = (pcrs) => {
  const infantIds = _.uniqWith(pcrs, (a, b) => a.infant.id === b.infant.id);
  const screens = infantIds.map((i) => compilePcrs(pcrs, i.infant.id));
  return screens;
};
