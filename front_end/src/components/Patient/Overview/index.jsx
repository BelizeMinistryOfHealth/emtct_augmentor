import React from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import { Box, Heading, Text } from 'grommet';
import Spinner from '../../Spinner';
import ErrorBoundary from '../../ErrorBoundary';
import PatientBasicInfo from '../PatientBasicInfo/PatientBasicInfo';
import PregnancyHistory from '../Pregnancy/PregnancyHistory';
import PatientIdSearch from '../../Search/PatientIdSearch';

const NoPatientFound = () => {
  return (
    <Box gap={'medium'} pad={'medium'} justify={'center'} align={'center'}>
      <Text size={'xlarge'}>No Patient Found!</Text>
    </Box>
  );
};

const PatientOverview = (props) => {
  const { patient } = props;
  return (
    <Box
      direction={'row-responsive'}
      gap={'small'}
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
  );
};

const Loading = () => {
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
};

const Overview = () => {
  const { patientId } = useParams();
  const history = useHistory();
  const { httpInstance } = useHttpApi();
  const [patient, setPatient] = React.useState();
  const [loading, setLoading] = React.useState(true);
  const [, setError] = React.useState(undefined);
  React.useEffect(() => {
    console.log({ patientId, loading });
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

  const patientIdSearchHandler = (pId) => {
    setLoading(true);
    history.push(`/patient/${pId}`, { id: pId });
  };

  return (
    <ErrorBoundary>
      <Box
        direction={'column'}
        fill={'horizontal'}
        gap={'large'}
        justify={'start'}
        align={'center'}
      >
        <Box align={'center'} pad={'large'}>
          <Box
            fill
            align={'center'}
            justify={'center'}
            direction={'row-responsive'}
          >
            <PatientIdSearch onSubmit={patientIdSearchHandler} />
          </Box>
        </Box>
        {loading && <Loading />}
        {!loading && patient && <PatientOverview patient={patient} />}
        {!loading && !patient && <NoPatientFound />}
      </Box>
    </ErrorBoundary>
  );
};

export default Overview;
