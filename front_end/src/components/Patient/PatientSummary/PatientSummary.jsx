import { Box, Heading, Text } from 'grommet';
import React from 'react';
import { useParams } from 'react-router-dom';
import { fetchPatient } from '../../../api/patient';
import { useHttpApi } from '../../../providers/HttpProvider';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import DiagnosisHistory from '../Diagnoses/Diagnoses';
import ObstetricHistory from '../ObstetricHistory/ObstetricHistory';
import PatientBasicInfo from '../PatientBasicInfo/PatientBasicInfo';
import Spinner from '../../Spinner';

const PatientSummary = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [patient, setPatient] = React.useState();
  const [loading, setLoading] = React.useState(true);
  const [, setError] = React.useState(undefined);
  React.useEffect(() => {
    const getPatient = () => {
      fetchPatient(patientId, httpInstance)
        .then((response) => {
          setPatient(response);
          setLoading(false);
        })
        .catch((e) => {
          setLoading(false);
          setPatient(undefined);
          setError(e.toJSON());
        });
    };
    if (loading) {
      getPatient();
    }
  }, [patientId, httpInstance, loading]);

  if (loading) {
    return (
      <Layout>
        <Box
          direction={'column'}
          fill={'horizontal'}
          gap={'large'}
          pad={'large'}
          justify={'center'}
          align={'center'}
        >
          <Heading>
            <Text>Loading </Text>
            <Spinner />
          </Heading>
        </Box>
      </Layout>
    );
  }

  if (patient && patient.basicInfo) {
    const {
      basicInfo,
      nextOfKin,
      obstetricHistory,
      diagnosesPrePregnancy,
    } = patient;
    return (
      <Layout location={props.location} {...props}>
        <ErrorBoundary>
          <Box
            direction={'column'}
            fill={'horizontal'}
            gap={'large'}
            justify={'center'}
          >
            <Box
              direction={'row-responsive'}
              gap={'medium'}
              pad={'medium'}
              justify={'start'}
              align={'start'}
            >
              <PatientBasicInfo basicInfo={basicInfo} nextOfKin={nextOfKin} />
              <ObstetricHistory obstetricHistory={obstetricHistory} />
            </Box>
            <Box gap={'medium'}>
              <DiagnosisHistory
                diagnosisHistory={diagnosesPrePregnancy}
                caption={'Illnesses before Pregnancy'}
              />
            </Box>
          </Box>
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
