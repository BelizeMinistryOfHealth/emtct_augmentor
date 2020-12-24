import {
  Box,
  Button,
  DateInput,
  Form,
  FormField,
  Heading,
  Text,
  TextArea,
  TextInput,
} from 'grommet';
import { FormPreviousLink } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import Layout from '../../Layout/Layout';

const ContraceptivesCreateForm = () => {
  const [contraceptive, setContraceptive] = React.useState();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const [patientData, setPatientData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });
  const { patientId, pregnancyId } = useParams();
  const history = useHistory();
  const { httpInstance } = useHttpApi();

  const onSubmit = (e) => {
    setContraceptive({ ...e.value, patientId: parseInt(patientId) });
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const fetchPatient = async () => {
      try {
        const result = await httpInstance.get(`/patients/${patientId}`);
        setPatientData({ data: result.data, loading: false, error: undefined });
      } catch (e) {
        console.error(e);
        setPatientData({ data: undefined, loading: false, error: e });
      }
    };
    if (patientData.loading) {
      fetchPatient();
    }
  }, [httpInstance, patientId, patientData]);

  React.useEffect(() => {
    const post = async (contraceptive) => {
      try {
        await httpInstance.post(
          `/patients/${patientId}/pregnancy/${pregnancyId}/contraceptivesUsed`,
          contraceptive
        );
        setStatus('SUCCESS');
      } catch (e) {
        console.error(e);
        setStatus('ERROR');
      }
    };
    if (status === 'SUBMIT') {
      post(contraceptive);
    }
  }, [
    contraceptive,
    httpInstance,
    status,
    pregnancyId,
    patientId,
    patientData,
  ]);

  if (status === 'SUCCESS') {
    return (
      <Layout>
        <Box
          fill={'vertical'}
          overflow={'auto'}
          pad={'medium'}
          width={'xlarge'}
          justify={'center'}
        >
          <Button
            icon={<FormPreviousLink size={'large'} />}
            onClick={() =>
              history.push(
                `/patient/${patientId}/pregnancy/${pregnancyId}/contraceptives`
              )
            }
          ></Button>
          <Box
            flex={false}
            direction={'row-responsive'}
            justify={'center'}
            fill={'horizontal'}
          >
            <Heading level={2} margin={'none'}>
              Successfully Saved Contraceptive Information!
            </Heading>
          </Box>
        </Box>
      </Layout>
    );
  }

  return (
    <Layout>
      <Box
        fill={'vertical'}
        overflow={'auto'}
        pad={'medium'}
        width={'xlarge'}
        justify={'center'}
      >
        <Button
          icon={<FormPreviousLink size={'large'} />}
          onClick={() =>
            history.push(
              `/patient/${patientId}/pregnancy/${pregnancyId}/contraceptives`
            )
          }
        ></Button>
        <Box
          direction={'column'}
          align={'start'}
          fill={'horizontal'}
          justify={'between'}
          alignContent={'center'}
        >
          <Text size={'xxlarge'} weight={'bold'} textAlign={'start'}>
            Create Contraceptive
          </Text>
          {patientData && patientData.data && (
            <Text size={'large'} textAlign={'end'} weight={'normal'}>
              {patientData.data.patient.firstName}{' '}
              {patientData.data.patient.lastName}
            </Text>
          )}
        </Box>
        {status === 'ERROR' && (
          <Box
            fill={'horizontal'}
            pad={'medium'}
            gap={'medium'}
            background={'red'}
          >
            <Text>Error creating contraceptive!</Text>
          </Box>
        )}

        <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
          <Form onSubmit={onSubmit}>
            <FormField label={'Contraceptive'} name={'contraceptive'} required>
              <TextInput placeholder={'Contraceptive'} name={'contraceptive'} />
            </FormField>
            <FormField label={'Comments'} name={'comments'} required>
              <TextArea name={'comments'} />
            </FormField>
            <FormField label={'Date used'} name={'dateUsed'} required>
              <DateInput format={'yyyy-mm-dd'} name={'dateUsed'} />
            </FormField>
            <Box flex={false} align={'center'}>
              <Button type={'submit'} label={'Save'} primary />
            </Box>
          </Form>
        </Box>
      </Box>
    </Layout>
  );
};

export default ContraceptivesCreateForm;
