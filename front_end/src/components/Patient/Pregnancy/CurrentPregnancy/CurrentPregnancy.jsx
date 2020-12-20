import { Box, Heading, Text } from 'grommet';
import React from 'react';
import { useParams } from 'react-router-dom';
import { fetchCurrentPregnancy } from '../../../../api/patient';
import { useHttpApi } from '../../../../providers/HttpProvider';
import Layout from '../../../Layout/Layout';
import PatientBasicInfo from '../../PatientBasicInfo/PatientBasicInfo';
import AppTabs from './Tabs/Tabs';
import PregnancyVitals from '../PregnancyVitals/PregnancyVitals';
import PreNatalCare from '../PreNatalCare/PreNatalCare';
import Spinner from '../../../Spinner';
import DiagnosisHistory from '../../Diagnoses/Diagnoses';

const BasicInfoComponent = ({ currentPregnancy }) => {
  return (
    <Box
      direction={'column'}
      pad={'medium'}
      gap={'large'}
      justify={'center'}
      fill={'horizontal'}
    >
      <Box
        direction={'row-responsive'}
        gap={'medium'}
        pad={'medium'}
        justify={'start'}
        align={'start'}
        fill={'horizontal'}
      >
        <PatientBasicInfo patient={currentPregnancy.patient} />
        <PregnancyVitals
          vitals={currentPregnancy.pregnancy}
          fill={'horizontal'}
        />
      </Box>
      <Box>
        <PreNatalCare info={currentPregnancy.pregnancy} />
      </Box>
    </Box>
  );
};

const PregnancyDiagnoses = ({ diagnoses }) => {
  return (
    <Box
      justify={'center'}
      align={'center'}
      fill={'horizontal'}
      gap={'large'}
      pad={'medium'}
    >
      <DiagnosisHistory
        diagnosisHistory={diagnoses}
        caption={'Illnesses during Pregnancy'}
      />
    </Box>
  );
};

const CurrentPregnancy = (props) => {
  const { location } = props;
  const { patientId, pregnancyId } = useParams();
  const { httpInstance } = useHttpApi();
  const [pregnancyData, setPregnancyData] = React.useState({
    currentPregnancy: {},
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getCurrentPregnancy = async () => {
      try {
        const result = await fetchCurrentPregnancy(
          patientId,
          pregnancyId,
          httpInstance
        );
        setPregnancyData({
          currentPregnancy: result,
          loading: false,
          error: undefined,
        });
      } catch (e) {
        // eslint-disable-next-line no-undef
        console.error(e);
        setPregnancyData({ currentPregnancy: {}, loading: false, error: e });
      }
    };
    if (pregnancyData.loading) {
      getCurrentPregnancy();
    }
  }, [httpInstance, pregnancyData, patientId, pregnancyId]);

  if (pregnancyData.loading) {
    return (
      <Layout>
        <Box
          direction={'column'}
          gap={'large'}
          pad={'large'}
          justify={'center'}
          align={'center'}
          fill
        >
          <Heading>
            <Text>Loading </Text>
            <Spinner />
          </Heading>
        </Box>
      </Layout>
    );
  }

  if (pregnancyData.error) {
    return (
      <Box
        direction={'column'}
        gap={'large'}
        pad={'large'}
        justify={'center'}
        align={'center'}
        fill
      >
        <Heading>
          <Text>Ooops. An error occurred while loading the data. </Text>
        </Heading>
      </Box>
    );
  }

  return (
    <Layout location={location} props={props}>
      <AppTabs
        basicInfo={
          <BasicInfoComponent
            currentPregnancy={pregnancyData.currentPregnancy}
          />
        }
        pregnancyDiagnoses={
          <PregnancyDiagnoses
            diagnoses={pregnancyData.currentPregnancy.diagnosesDuringPregnancy}
          />
        }
        diagnosesBeforePregnancy={
          <PregnancyDiagnoses
            diagnoses={pregnancyData.currentPregnancy.diagnosesBeforePregnancy}
          />
        }
      />
    </Layout>
  );
};

export default CurrentPregnancy;
