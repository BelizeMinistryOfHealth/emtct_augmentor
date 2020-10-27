import { Box, Text } from 'grommet';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useRecoilValueLoadable } from 'recoil';
import { patientSelector } from '../../../state';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import DiagnosisHistory from '../Diagnoses/Diagnoses';
import ObstetricHistory from '../ObstetricHistory/ObstetricHistory';
import PatientBasicInfo from '../PatientBasicInfo/PatientBasicInfo';

const PatientSummary = (props) => {
  const { patientId } = useParams();

  const { state, contents } = useRecoilValueLoadable(
    patientSelector(patientId)
  );
  let patient = {};
  switch (state) {
    case 'hasValue':
      patient = contents;
      break;
    case 'hasError':
      console.dir({ contents });
      return contents.message;
    case 'loading':
      return 'Loading....';
    default:
      return '';
  }

  if (patient) {
    const {
      basicInfo,
      nextOfKin,
      obstetricHistory,
      diagnosesPrePregnancy,
    } = patient;
    return (
      <Layout location={props.location} {...props}>
        <ErrorBoundary>
          <PatientBasicInfo basicInfo={basicInfo} nextOfKin={nextOfKin} />
          <ObstetricHistory obstetricHistory={obstetricHistory} />
          <DiagnosisHistory
            diagnosisHistory={diagnosesPrePregnancy}
            caption={'Illnesses before Pregnancy'}
          />
        </ErrorBoundary>
      </Layout>
    );
  }
  return (
    <Box gap={'medium'} pad={'medium'} justify={'center'} align={'center'}>
      <Text size={'xlarge'}>No Patient Found!</Text>
    </Box>
  );
};

export default PatientSummary;
