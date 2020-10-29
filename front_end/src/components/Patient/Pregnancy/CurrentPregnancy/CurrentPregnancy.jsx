import { Box, Text } from 'grommet';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useRecoilValueLoadable } from 'recoil';
import {
  currentPregnancySelector,
  pregnancyLabResultsSelector,
} from '../../../../state';
import Layout from '../../../Layout/Layout';
import ArvTreatment from '../../ArvTreatment/ArvTreatment';
import DiagnosisHistory from '../../Diagnoses/Diagnoses';
import PatientBasicInfo from '../../PatientBasicInfo/PatientBasicInfo';
import AppTabs from '../../Tabs/Tabs';
import LabResults from '../LabResults/LabResults';
import PregnancyVitals from '../PregnancyVitals/PregnancyVitals';
import PreNatalCare from '../PreNatalCare/PreNatalCare';

const BasicInfoComponent = ({ currentPregnancy }) => {
  return (
    <Box
      direction={'row-responsive'}
      gap={'medium'}
      pad={'medium'}
      justify={'start'}
      align={'start'}
      fill
    >
      <PatientBasicInfo
        basicInfo={currentPregnancy.basicInfo}
        nextOfKin={currentPregnancy.nextOfKin}
      />
      <PregnancyVitals vitals={currentPregnancy.vitals} />
    </Box>
  );
};

const Arvs = ({ currentPregnancy }) => {
  return (
    <Box
      gap={'medium'}
      pad={'medium'}
      justify={'center'}
      align={'center'}
      direction={'row-responsive'}
      fill
    >
      <PreNatalCare info={currentPregnancy.prenatalCareInfo} />
      <ArvTreatment
        patientId={currentPregnancy.basicInfo.patientId}
        encounterId={currentPregnancy.encounterId}
      />
      <DiagnosisHistory
        diagnosisHistory={currentPregnancy.pregnancyDiagnoses}
        caption={'Illnesses during Pregnancy'}
      />
    </Box>
  );
};

const LabTests = ({ patientId, encounterId }) => {
  const { state, contents } = useRecoilValueLoadable(
    pregnancyLabResultsSelector(patientId)
  );

  switch (state) {
    case 'hasValue':
      return (
        <Box
          gap={'medium'}
          pad={'medium'}
          justify={'center'}
          align={'center'}
          direction={'row-responsive'}
        >
          <LabResults
            labResults={contents}
            caption={
              <Text weight={'bold'}>Lab Test Results During Pregnancy</Text>
            }
          />
        </Box>
      );
    case 'hasError':
      return contents.message;
    case 'loading':
      return 'loading';
    default:
      return '';
  }
};

const CurrentPregnancy = (props) => {
  const { location } = props;
  const { id } = useParams();

  const { state, contents } = useRecoilValueLoadable(
    currentPregnancySelector(id)
  );
  let currentPregnancy = {};
  switch (state) {
    case 'hasValue':
      currentPregnancy = contents;
      break;
    case 'hasError':
      return contents.message;
    case 'loading':
      return 'loading';
    default:
      return '';
  }

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
        <AppTabs
          basicInfo={<BasicInfoComponent currentPregnancy={currentPregnancy} />}
          arvs={<Arvs currentPregnancy={currentPregnancy} />}
          labResults={
            <LabTests
              patientId={id}
              encounterId={currentPregnancy.encounterId}
            />
          }
        />
      </Box>
    </Layout>
  );
};

export default CurrentPregnancy;
