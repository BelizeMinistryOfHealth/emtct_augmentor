import { atom, selectorFamily } from 'recoil';
import {
  fetchPatient,
  fetchCurrentPregnancy,
  fetchArvsTreatment,
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
    return await fetchPatient(patientId);
  },
});

export const currentPregnancyAtom = atom({
  key: 'currentPregnancy',
  default: {},
});

export const currentPregnancySelector = selectorFamily({
  key: 'getCurrentPregnancyAPI',
  get: (patientId) => async () => {
    return await fetchCurrentPregnancy(patientId);
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
