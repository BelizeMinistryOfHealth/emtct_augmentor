import { format, parseISO } from 'date-fns';
import { Box, Heading, Text } from 'grommet';
import { InProgress } from 'grommet-icons';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useHttpApi } from '../../providers/HttpProvider';
import AppCard from '../AppCard/AppCard';
import Layout from '../Layout/Layout';
import InfantTabs from './InfantTabs';

const BasicInfoHeaders = () => {
  return (
    <Box pad={'large'} gap={'large'}>
      <Text size={'large'} weight={'bold'} textAlign={'start'}>
        Name:
      </Text>
      <Text size={'large'} weight={'bold'} textAlign={'start'}>
        Patient Id:
      </Text>
      <Text size={'large'} textAlign={'start'} weight={'bold'}>
        Date of Birth:
      </Text>
      <Text size={'large'} weight={'bold'} textAlign={'start'}>
        Mother:
      </Text>
    </Box>
  );
};

const BasicInfo = ({ data }) => {
  return (
    <Box gap={'medium'} align={'center'} fill={'horizontal'}>
      <Box direction={'row'} gap={'medium'} fill={'horizontal'}>
        <BasicInfoHeaders />
        <Box pad={'large'} gap={'large'}>
          <Text size={'large'} textAlign={'start'}>
            {data.firstName} {data.lastName}
          </Text>
          <Text size={'large'} textAlign={'start'}>
            {data.id}
          </Text>
          <Text size={'large'} textAlign={'start'}>
            {format(parseISO(data.dob), 'dd LLL yyyy')}
          </Text>
          <Text size={'large'} textAlign={'start'}>
            {data.mother.firstName} {data.mother.lastName}
          </Text>
        </Box>
      </Box>
    </Box>
  );
};

const InfantInfo = ({ data }) => {
  return (
    <AppCard fill={'horizontal'}>
      <Box
        margin={'small'}
        direction={'column'}
        gap={'small'}
        alignSelf={'start'}
      >
        <BasicInfo data={data} />
      </Box>
    </AppCard>
  );
};

const Infant = () => {
  const { patientId, pregnancyId } = useParams();
  const { httpInstance } = useHttpApi();
  const [infantData, setInfantData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getInfant = () => {
      httpInstance
        .get(`/patients/${patientId}/pregnancy/${pregnancyId}/infant`)
        .then((result) => {
          // eslint-disable-next-line no-undef
          console.dir({ result });
          setInfantData({
            data: result.data,
            loading: false,
            error: undefined,
          });
        })
        .catch((e) => {
          if (e.response.status === 404) {
            setInfantData({
              data: undefined,
              loading: false,
              error: 'NOT_FOUND',
            });
          } else {
            // eslint-disable-next-line no-undef
            console.error({ error: e.toJSON(), status: e.response.status });
            setInfantData({
              data: undefined,
              loading: false,
              error: e,
            });
          }
        });
    };

    if (infantData.loading) {
      getInfant();
    }
  }, [infantData, httpInstance, patientId, pregnancyId]);

  if (infantData.loading) {
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
            <Text>Loading...</Text>
            <InProgress />
          </Heading>
        </Box>
      </Layout>
    );
  }

  if (infantData.error && infantData.error === 'NOT_FOUND') {
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
            <Text>No infant found four current pregnancy.</Text>
          </Heading>
        </Box>
      </Layout>
    );
  }

  if (infantData.error) {
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
            <Text>Ooops. An error occurred while loading the data.</Text>
          </Heading>
        </Box>
      </Layout>
    );
  }

  return (
    <Layout>
      <Box
        direction={'column'}
        pad={{ left: 'small', bottom: 'xxsmall' }}
        alignContent={'start'}
        alignSelf={'start'}
        fill={'horizontal'}
      >
        <InfantTabs data={infantData.data} pregnancyId={pregnancyId}>
          <InfantInfo data={infantData.data} />
        </InfantTabs>
      </Box>
    </Layout>
  );
};

export default Infant;

export * from './HivScreenings';
