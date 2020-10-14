export const fetchPatient = async (patientId) => {
  const basicInfo = {
    firstName: 'Jane',
    lastName: 'Doe',
    dob: '2000-10-21',
    ssn: '145134235',
    patientId,
    countryOfBirth: 'Belize',
    district: 'Cayo',
    community: 'Belmopan',
    address: 'Corozal Street',
    education: 'High School',
    ethnicity: 'Ethnic Group',
    hiv: false,
  };
  const nextOfKin = {
    name: 'John Doe',
    phoneNumber: '6632888',
  };

  const obstetricHistory = [
    { id: 1, date: '2010-02-21', event: 'Live Born' },
    { id: 2, date: '2012-11-30', event: 'Miscarriage' },
    { id: 3, date: '2014-10-01', event: 'Live Born' },
  ];

  const diagnosesPrePregnancy = [
    { id: 1, date: '2008-10-21', name: 'common cold' },
    { id: 2, date: '2009-04-10', name: 'rash' },
    { id: 3, date: '2009-09-21', name: 'common cold' },
    { id: 4, date: '2010-02-23', name: 'conjuctivitis' },
  ];

  return new Promise((resolve) => {
    return resolve({
      basicInfo,
      nextOfKin,
      obstetricHistory,
      diagnosesPrePregnancy,
    });
  });
};

export const fetchCurrentPregnancy = async (patientId) => {
  const basicInfo = {
    firstName: 'Jane',
    lastName: 'Doe',
    dob: '2000-10-21',
    ssn: '145134235',
    patientId,
    countryOfBirth: 'Belize',
    district: 'Cayo',
    community: 'Belmopan',
    address: 'Corozal Street',
    education: 'High School',
    ethnicity: 'Ethnic Group',
    hiv: false,
  };

  const nextOfKin = {
    name: 'John Doe',
    phoneNumber: '6632888',
  };
  const prenatalCareInfo = {
    dateOfBooking: '2020-07-01',
    gestationAge: 7,
    prenatalCareProvider: 'Public',
    totalChecks: 4,
  };

  const currentPregnancy = {
    encounterId: 2121,
    vitals: {
      gestationalAge: 4,
      para: 10,
      cs: false,
      abortiveOutcome: 'None',
      diagnosisDate: '2020-06-15',
      planned: false,
      ageAtLmp: 19,
      LMP: '2020-04-21',
      EDD: '2021-01-21',
    },
    basicInfo,
    nextOfKin,
    prenatalCareInfo,
  };

  return new Promise((resolve) => {
    resolve(currentPregnancy);
  });
};
