import { Box } from 'grommet';
import React from 'react';
import PatientBasicInfo from './PatientBasicInfo';

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

  return (
    <Box
      align={'center'}
      justify={'start'}
      direction={'row-responsive'}
      gap={'medium'}
      pad={'medium'}
      fill
    >
      <PatientBasicInfo basicInfo={basicInfo} nextOfKin={nextOfKin} />
    </Box>
  );
};

export default PatientSummary;
