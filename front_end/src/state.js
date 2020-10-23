import { useHttpApi } from './providers/HttpProvider';
import { atom, selectorFamily } from 'recoil';
import {
  fetchPatient,
  fetchCurrentPregnancy,
  fetchArvsTreatment,
  fetchPregnancyLabResults,
} from './api/patient';

export const patientIdState = atom({
  key: 'patientId',
  default: '',
});

export const patientAtom = atom({
  key: 'patient',
  default: {},
});

export const patientSelector = selectorFamily({
  key: 'getPatientAPI',
  get: (patientId) => async () => {
    const { httpInstance } = useHttpApi();
    return await fetchPatient(patientId, httpInstance);
  },
});

export const currentPregnancyAtom = atom({
  key: 'currentPregnancy',
  default: {},
});

export const currentPregnancySelector = selectorFamily({
  key: 'getCurrentPregnancyAPI',
  get: (patientId) => async () => {
    const { httpInstance } = useHttpApi();
    return await fetchCurrentPregnancy(patientId, httpInstance);
  },
});

export const arvTreatments = atom({
  key: 'arvTreatments',
  default: [],
});

export const arvTreatmentsSelector = selectorFamily({
  key: 'arvTreatmentsAPI',
  get: (patientId, encounterId) => async () => {
    return await fetchArvsTreatment(patientId, encounterId);
  },
});

export const pregnancyLabResultsAtom = atom({
  key: 'pregnancyLabResults',
  default: [],
});

export const pregnancyLabResultsSelector = selectorFamily({
  key: 'pregnancyLabResultsAPI',
  get: (patientId) => async () => {
    const { httpInstance } = useHttpApi();
    return await fetchPregnancyLabResults(patientId, httpInstance);
  },
});
