import { Box, Card, CardBody } from 'grommet';
import React from 'react';
import ObstetricHistory from '../ObstetricHistory/ObstetricHistory';
import PatientBasicInfo from '../PatientBasicInfo/PatientBasicInfo';

const PatientSummary = (props) => {
  const id = props.location.state.id;
  // const { generalInformation, currentPregnancy } = patient;
  const basicInfo = {
    firstName: 'Jane',
    lastName: 'Doe',
    dob: '2000-10-21',
    ssn: '145134235',
    patientId: id,
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

  return (
    <Box
      align={'start'}
      justify={'start'}
      direction={'row-responsive'}
      gap={'medium'}
      pad={'medium'}
      fill
    >
      <PatientBasicInfo basicInfo={basicInfo} nextOfKin={nextOfKin} />
      <ObstetricHistory obstetricHistory={obstetricHistory} />
    </Box>
  );
};

export default PatientSummary;
