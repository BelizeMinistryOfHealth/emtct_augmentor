import { Box, Nav } from 'grommet';
import React from 'react';
import { useRecoilState } from 'recoil';
import { useRecoilValue } from 'recoil';
import { useRecoilValueLoadable } from 'recoil';
import { patientSelector } from '../../../state';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import DiagnosisHistory from '../Diagnoses/Diagnoses';
import ObstetricHistory from '../ObstetricHistory/ObstetricHistory';
import PatientBasicInfo from '../PatientBasicInfo/PatientBasicInfo';

const PatientSummary = (props) => {
  const id = props.location.state.id;

  const { state, contents } = useRecoilValueLoadable(patientSelector(id));
  let patient = {};
  switch (state) {
    case 'hasValue':
      patient = contents;
      break;
    case 'hasError':
      return contents.message;
    case 'loading':
      return 'Loading....';
    default:
      return '';
  }
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
};

export default PatientSummary;
