import _ from 'lodash';
import React from 'react';
import { useHttpApi } from './http';

export const fetchPatient = async (patientId, httpInstance) => {
  const result = await httpInstance.get(`/patient/${patientId}`);
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

export const fetchCurrentPregnancy = async (patientId, httpInstance) => {
  const patientsData = await httpInstance.get(`/patient/${patientId}`);
  if (!patientsData.data) {
    return null;
  }
  const patient = patientsData.data;
  const basicInfo = patient.patient;
  const nextOfKin = {
    name: basicInfo.nextOfKin,
    phoneNumber: basicInfo.nextOfKinPhone,
  };

  const pregnancyData = await httpInstance.get(
    `/patient/${patientId}/currentPregnancy`
  );
  if (!pregnancyData.data) {
    return null;
  }

  const vitals = pregnancyData.data.vitals;
  if (_.isEmpty(vitals.abortiveOutcome.trim())) {
    vitals.abortiveOutcome = 'N/A';
  }

  const prenatalCareInfo = {
    dateOfBooking: vitals.dateOfBooking,
    gestationAge: vitals.gestationalAge,
    prenatalCareProvider: vitals.prenatalCareProvider,
    totalChecks: vitals.totalChecks,
  };

  const pregnancyDiagnoses = pregnancyData.data.diagnoses ?? [];

  return {
    vitals,
    basicInfo,
    nextOfKin,
    prenatalCareInfo,
    pregnancyDiagnoses,
  };
};

export const fetchHomeVisits = async (patientId, httpInstance) => {
  const result = await httpInstance.get(`/patient/${patientId}/homeVisits`);
  if (!result.data) {
    return [];
  }

  return result.data;
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
  }, [httpInstance]);
  return [data, loading, error];
};

export const fetchArvsTreatment = (patientId, encounterId) => {
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

export const fetchPregnancyLabResults = async (patientId, httpInstance) => {
  const labResultsResponse = await httpInstance.get(
    `/patient/${patientId}/currentPregnancy/labResults`
  );

  return labResultsResponse.data;
};
