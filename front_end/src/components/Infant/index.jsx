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
            {data.infant.firstName} {data.infant.lastName}
          </Text>
          <Text size={'large'} textAlign={'start'}>
            {data.infant.patientId}
          </Text>
          <Text size={'large'} textAlign={'start'}>
            {format(parseISO(data.infant.dob), 'dd LLL yyyy')}
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
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [infantData, setInfantData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getInfant = async () => {
      try {
        const result = await httpInstance.get(`/patient/${patientId}/infant`);
        setInfantData({
          data: result.data,
          loading: false,
          error: undefined,
        });
      } catch (e) {
        console.error(e);
        setInfantData({
          data: undefined,
          loading: false,
          error: e,
        });
      }
    };

    if (infantData.loading) {
      getInfant();
    }
  }, [infantData, httpInstance, patientId]);

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

  if (infantData.error) {
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
    </Layout>;
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
        <InfantTabs
          content={<InfantInfo data={infantData.data} />}
          data={infantData.data}
        />
      </Box>
    </Layout>
  );
};

export default Infant;

export * from './HivScreenings';
