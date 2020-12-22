import { Box, CardBody, Text } from 'grommet';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import AppCardHeader from '../../AppCard/AppCardHeader';
import Layout from '../../Layout/Layout';
import Prescriptions from '../../Prescriptions';
import Spinner from '../../Spinner';
import InfantTabs from '../InfantTabs';

const InfantSyphillisTreatment = (props) => {
  const { patientId, infantId, pregnancyId } = useParams();
  const { httpInstance } = useHttpApi();
  const [treatmentData, setTreatmentData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getTreatment = async () => {
      try {
        const result = await httpInstance.get(
          `/patients/${patientId}/infant/${infantId}/syphilisTreatments`
        );
        setTreatmentData({
          data: result.data,
          loading: false,
          error: undefined,
        });
      } catch (e) {
        // eslint-disable-next-line no-undef
        console.error(e);
        setTreatmentData({ data: undefined, loading: false, error: e });
      }
    };
    if (treatmentData.loading) {
      getTreatment();
    }
  }, [infantId, patientId, treatmentData, httpInstance]);

  if (treatmentData.loading) {
    return (
      <Layout props={props}>
        <Box
          direction={'column'}
          gap={'medium'}
          pad={'medium'}
          justify={'evenly'}
          align={'center'}
          fill
        >
          <Text>Loading....</Text>
          <Spinner />
        </Box>
      </Layout>
    );
  }

  if (treatmentData.error) {
    return (
      <Layout>
        <Box
          direction={'column'}
          gap={'medium'}
          pad={'medium'}
          justify={'evenly'}
          align={'center'}
          fill
        >
          <Text>Ooops.. An error occurred.</Text>
        </Box>
      </Layout>
    );
  }

  return (
    <Layout>
      <Box
        direction={'column'}
        gap={'medium'}
        pad={'medium'}
        justify={'evenly'}
        align={'center'}
        fill
      >
        <InfantTabs data={treatmentData.data.infant} pregnancyId={pregnancyId}>
          <AppCard fill={'horizontal'}>
            <AppCardHeader
              gap={'medium'}
              pad={'medium'}
              title={'Syphilis Treatments'}
              patient={treatmentData.data.infant}
            />
            <CardBody gap={'medium'} pad={'medium'}>
              <Prescriptions data={treatmentData.data} />
            </CardBody>
          </AppCard>
        </InfantTabs>
      </Box>
    </Layout>
  );
};

export default InfantSyphillisTreatment;
