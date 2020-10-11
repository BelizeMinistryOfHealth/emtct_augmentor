import { Box, Nav } from 'grommet';
import React from 'react';
import Layout from '../../Layout/Layout';
import DiagnosisHistory from '../Diagnoses/Diagnoses';
import ObstetricHistory from '../ObstetricHistory/ObstetricHistory';
import PatientBasicInfo from '../PatientBasicInfo/PatientBasicInfo';

const PatientSummary = (props) => {
  const id = props.location.state.id;
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

  const diagnosesPrePregnancy = [
    { id: 1, date: '2008-10-21', name: 'common cold' },
    { id: 2, date: '2009-04-10', name: 'rash' },
    { id: 3, date: '2009-09-21', name: 'common cold' },
    { id: 4, date: '2010-02-23', name: 'conjuctivitis' },
  ];

  return (
    <Layout>
      <PatientBasicInfo basicInfo={basicInfo} nextOfKin={nextOfKin} />
      <ObstetricHistory obstetricHistory={obstetricHistory} />
      <DiagnosisHistory
        diagnosisHistory={diagnosesPrePregnancy}
        caption={'Illnesses before Pregnancy'}
      />
    </Layout>
  );
};

export default PatientSummary;
