import React from 'react';
import { useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import { Box, Heading, Text } from 'grommet';
import Spinner from '../../Spinner';
import ErrorBoundary from '../../ErrorBoundary';
import PatientBasicInfo from '../PatientBasicInfo/PatientBasicInfo';
import PregnancyHistory from '../Pregnancy/PregnancyHistory';

const Overview = () => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [patient, setPatient] = React.useState();
  const [loading, setLoading] = React.useState(true);
  const [, setError] = React.useState(undefined);
  React.useEffect(() => {
    const getPatient = () => {
      httpInstance
        .get(`/patients/${patientId}`)
        .then((response) => {
          const pt = response.data.patient;
          pt.pregnancies = response.data.pregnancies;
          setPatient(pt);
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
      <Box
        direction={'row'}
        fill={'horizontal'}
        gap={'large'}
        justify={'start'}
        align={'start'}
      >
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
      </Box>
    );
  }

  if (patient) {
    return (
      <ErrorBoundary>
        <Box
          direction={'row'}
          fill={'horizontal'}
          gap={'large'}
          justify={'start'}
          align={'start'}
        >
          <Box
            direction={'row-responsive'}
            gap={'medium'}
            pad={'medium'}
            justify={'start'}
            align={'center'}
          >
            <PatientBasicInfo patient={patient} />
          </Box>
          <Box gap={'medium'} pad={'medium'}>
            <PregnancyHistory pregnancies={patient.pregnancies} />
          </Box>
        </Box>
      </ErrorBoundary>
    );
  }
  return (
    <Box gap={'medium'} pad={'medium'} justify={'center'} align={'center'}>
      <Text size={'xlarge'}>No Patient Found!</Text>
    </Box>
  );
};

export default Overview;
