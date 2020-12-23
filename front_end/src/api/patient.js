import _ from 'lodash';
import React from 'react';
import { useHttpApi } from './http';

export const fetchPatient = async (patientId, httpInstance) => {
  const result = await httpInstance.get(`/patients/${patientId}`);
  const data = result.data;

  if (!data.patient) {
    return null;
  }
  const basicInfo = data.patient;
  const nextOfKin = {
    name: basicInfo.nextOfKin,
    phoneNumber: basicInfo.nextOfKinPhone,
  };

  const obstetricHistory = data.obstetricHistory;

  const diagnosesPrePregnancy = data.diagnoses;

  return {
    basicInfo,
    nextOfKin,
    obstetricHistory,
    diagnosesPrePregnancy,
  };
};

export const fetchCurrentPregnancy = async (
  patientId,
  pregnancyId,
  httpInstance
) => {
  const pregnancyData = await httpInstance.get(
    `/patients/${patientId}/pregnancy/${pregnancyId}`
  );
  const data = pregnancyData.data;
  const pregnancy = data.pregnancy;
  const interval = data.interval;
  pregnancy.interval = interval;

  const pregnancyOutcome = pregnancy.obstetricDetails.pregnancyOutcome;
  if (!pregnancyOutcome || _.isEmpty(pregnancyOutcome.trim())) {
    pregnancy.obstetricDetails.pregnancyOutcome = pregnancyOutcome;
  }

  let gestationalAge = 'N/A';
  if (pregnancy.anc.gestationalAge > 7) {
    gestationalAge = `${Math.ceil(pregnancy.anc.gestationalAge / 7)} weeks`;
  }

  if (pregnancy.anc.gestationalAge < 7) {
    gestationalAge = `${pregnancy.anc.gestationalAge} days`;
  }
  pregnancy.anc.gestationalAge = gestationalAge;

  const pregnancyDiagnoses = pregnancyData.data.diagnoses ?? [];

  return {
    pregnancy,
    pregnancyDiagnoses,
    patient: data.patient,
    diagnosesDuringPregnancy: data.diagnosesDuringPregnancy,
    diagnosesBeforePregnancy: data.diagnosesBeforePregnancy,
  };
};

export const useEditHomeVisit = async (homeVisit) => {
  const { httpInstance } = useHttpApi();
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState(undefined);
  const [data, setData] = React.useState(homeVisit);

  React.useEffect(() => {
    const edit = async () => {
      setLoading(true);
      try {
        const result = await httpInstance.put(
          `/patient/homeVisit/${homeVisit.id}`,
          homeVisit
        );
        setData(result.data);
        setLoading(false);
      } catch (e) {
        console.error(e);
        setError('Edit request failed');
        setLoading(false);
      }
    };
    edit();
  }, [httpInstance, homeVisit]);
  return [data, loading, error];
};

export const fetchArvsTreatment = (patientId, encounterId) => {
  console.log({ patientId, encounterId });
  const treatments = [
    {
      id: 1,
      arvName: 'ARV 1',
      dosage: '2 every 6 hours',
      date: '2019-01-21',
      comments: '',
    },
    {
      id: 2,
      arvName: 'ARV 1',
      dosage: '2 every 6 hours',
      date: '2019-06-30',
      comments: '',
    },
    {
      id: 3,
      arvName: 'ARV 2',
      dosage: '2 every 6 hours',
      date: '2020-02-10',
      comments: '',
    },
  ];

  return new Promise((resolve) => {
    resolve(treatments);
  });
};

export const fetchPregnancyLabResults = async (
  patientId,
  pregnancyId,
  httpInstance
) => {
  const labResultsResponse = await httpInstance.get(
    `/patients/${patientId}/pregnancy/${pregnancyId}/labResults`
  );

  return labResultsResponse.data;
};
