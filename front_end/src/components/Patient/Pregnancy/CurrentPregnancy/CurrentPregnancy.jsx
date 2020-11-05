import { Box, Heading, Text } from 'grommet';
import { InProgress } from 'grommet-icons';
import React from 'react';
import { useParams } from 'react-router-dom';
import { fetchCurrentPregnancy } from '../../../../api/patient';
import { useHttpApi } from '../../../../providers/HttpProvider';
import Layout from '../../../Layout/Layout';
import ArvTreatment from '../../ArvTreatment/ArvTreatment';
import DiagnosisHistory from '../../Diagnoses/Diagnoses';
import PatientBasicInfo from '../../PatientBasicInfo/PatientBasicInfo';
import AppTabs from '../../Tabs/Tabs';
import PregnancyVitals from '../PregnancyVitals/PregnancyVitals';
import PreNatalCare from '../PreNatalCare/PreNatalCare';

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
        fill
      >
        <PatientBasicInfo
          basicInfo={currentPregnancy.basicInfo}
          nextOfKin={currentPregnancy.nextOfKin}
        />
        <PregnancyVitals vitals={currentPregnancy.vitals} />
      </Box>
      <Box>
        <PreNatalCare info={currentPregnancy.prenatalCareInfo} />
      </Box>
    </Box>
  );
};

const PregnancyDiagnoses = ({ currentPregnancy }) => {
  return (
    <Box
      justify={'center'}
      align={'center'}
      fill={'horizontal'}
      gap={'large'}
      pad={'medium'}
    >
      <DiagnosisHistory
        diagnosisHistory={currentPregnancy.pregnancyDiagnoses}
        caption={'Illnesses during Pregnancy'}
      />
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
      <ArvTreatment
        patientId={currentPregnancy.basicInfo.patientId}
        encounterId={currentPregnancy.encounterId}
      />
    </Box>
  );
};

const CurrentPregnancy = (props) => {
  const { location } = props;
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [pregnancyData, setPregnancyData] = React.useState({
    currentPregnancy: {},
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getCurrentPregnancy = async () => {
      try {
        const result = await fetchCurrentPregnancy(patientId, httpInstance);
        setPregnancyData({
          currentPregnancy: result,
          loading: false,
          error: undefined,
        });
      } catch (e) {
        console.error(e);
        setPregnancyData({ currentPregnancy: {}, loading: false, error: e });
      }
    };
    if (pregnancyData.loading) {
      getCurrentPregnancy();
    }
  }, [httpInstance, pregnancyData, patientId, setPregnancyData]);

  if (pregnancyData.loading) {
    return (
      <Box
        direction={'colomn'}
        gap={'large'}
        pad={'large'}
        justify={'center'}
        align={'center'}
        fill
      >
        <Heading>
          <Text>Loading </Text>
          <InProgress />
        </Heading>
      </Box>
    );
  }

  if (pregnancyData.error) {
    return (
      <Box
        direction={'colomn'}
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
      <Box
        direction={'column'}
        gap={'medium'}
        pad={'medium'}
        justify={'start'}
        align={'start'}
        fill
      >
        <AppTabs
          basicInfo={
            <BasicInfoComponent
              currentPregnancy={pregnancyData.currentPregnancy}
            />
          }
          arvs={<Arvs currentPregnancy={pregnancyData.currentPregnancy} />}
          diagnoses={
            <PregnancyDiagnoses
              currentPregnancy={pregnancyData.currentPregnancy}
            />
          }
        />
      </Box>
    </Layout>
  );
};

export default CurrentPregnancy;
