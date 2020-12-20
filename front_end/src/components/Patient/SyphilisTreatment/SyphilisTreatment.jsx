import { Box, CardBody, Heading, Text } from 'grommet';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import AppCardHeader from '../../AppCard/AppCardHeader';
import Layout from '../../Layout/Layout';
import Prescriptions from '../../Prescriptions';

const SyphilisTreatment = (props) => {
  const { patientId, pregnancyId } = useParams();
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
          `/patients/${patientId}/pregnancy/${pregnancyId}/syphilisTreatments`
        );
        if (result.status === 204) {
          setTreatmentData({
            data: undefined,
            loading: false,
            error: new Error('no data found'),
          });
        } else {
          setTreatmentData({
            data: result.data,
            loading: false,
            error: undefined,
          });
        }
      } catch (e) {
        // eslint-disable-next-line no-undef
        console.error(e);
        setTreatmentData({ data: undefined, loading: false, error: e });
      }
    };
    if (treatmentData.loading) {
      getTreatment();
    }
  }, [treatmentData, httpInstance, patientId, pregnancyId]);

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
          <Text>Loading...</Text>
        </Box>
      </Layout>
    );
  }

  if (treatmentData.error) {
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
          <Heading level={2}>No Data Found</Heading>
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
        <AppCard fill={'horizontal'}>
          <AppCardHeader
            gap={'medium'}
            pad={'medium'}
            title={'Syphilis Treatments'}
            patient={treatmentData.data.patient}
          />
          <CardBody gap={'medium'} pad={'medium'}>
            <Prescriptions data={treatmentData.data}></Prescriptions>
          </CardBody>
        </AppCard>
      </Box>
    </Layout>
  );
};

export default SyphilisTreatment;
