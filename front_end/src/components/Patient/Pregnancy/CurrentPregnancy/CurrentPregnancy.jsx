import { Box } from 'grommet';
import React from 'react';
import Layout from '../../../Layout/Layout';
import PatientBasicInfo from '../../PatientBasicInfo/PatientBasicInfo';
import PregnancyVitals from '../PregnancyVitals/PregnancyVitals';
import PreNatalCare from '../PreNatalCare/PreNatalCare';

const CurrentPregnancy = (props) => {
  const { location } = props;

  const id = location.state.id;
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
  };

  const prenatalCareInfo = {
    dateOfBooking: '2020-07-01',
    gestationAge: 7,
    prenatalCareProvider: 'Public',
  };

  return (
    <Layout location={location} props={props}>
      <Box
        direction={'column'}
        gap={'medium'}
        pad={'medium'}
        justify={'start'}
        align={'start'}
        fill
      >
        <PatientBasicInfo basicInfo={basicInfo} nextOfKin={nextOfKin} />
        <Box
          direction={'row-responsive'}
          gap={'medium'}
          pad={'medium'}
          justify={'start'}
          align={'start'}
          fill
        >
          <PregnancyVitals vitals={currentPregnancy.vitals} />
          <PreNatalCare info={prenatalCareInfo} />
        </Box>
      </Box>
    </Layout>
  );
};

export default CurrentPregnancy;
