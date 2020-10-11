import axios from 'axios';
// Given an encounter id, return the patient's basic demographic information and current pregnancy information.

// If a result is found we change routes.
const baseUrl = '';
export const useGetPatientById = async (patientId) => {
  // We need some http client here at some point. Maybe it should come from props. Oh, this should be a data provider.
  try {
    const results = await axios.get(`${baseUrl}/patientId`);
    return results;
  } catch (e) {
    console.error('oh oh, an error occurred fetching patient by id: ', e);
    return null;
  }
};
