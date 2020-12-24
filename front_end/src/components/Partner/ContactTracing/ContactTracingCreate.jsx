import {
  Box,
  Button,
  DateInput,
  Form,
  FormField,
  Heading,
  TextArea,
  TextInput,
} from 'grommet';
import { FormPreviousLink } from 'grommet-icons';
import React from 'react';
import { Redirect, useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import Layout from '../../Layout/Layout';
import Spinner from '../../Spinner';

const CcontactTracingCreate = () => {
  const [patientData, setPatientData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const [contactTracing, setContactTracing] = React.useState();
  const { patientId, pregnancyId } = useParams();
  const { httpInstance } = useHttpApi();
  const history = useHistory();

  const onSubmit = (e) => {
    e.preventDefault();
    setContactTracing({ ...e.value, patientId: parseInt(patientId) });
    setStatus('SUBMIT');
  };

  // Retrieve patient information
  React.useEffect(() => {
    const fetchPatient = () => {
      httpInstance
        .get(`/patients/${patientId}`)
        .then((response) => {
          setPatientData({
            data: response.data,
            loading: false,
            error: undefined,
          });
        })
        .catch((e) => {
          console.error(e);
          setPatientData({
            data: undefined,
            loading: false,
            error: e.toJSON(),
          });
        });
    };
    if (patientData.loading) {
      fetchPatient();
    }
  }, [patientId, httpInstance, patientData]);

  React.useEffect(() => {
    const post = () => {
      httpInstance
        .post(
          `/patients/${patientId}/pregnancy/${pregnancyId}/contactTracing`,
          contactTracing
        )
        .then(() => {
          setStatus('SUCCESS');
        })
        .catch((e) => {
          console.error(e);
          setStatus('ERROR');
        });
    };
    if (status === 'SUBMIT') {
      post();
    }
  }, [patientId, pregnancyId, httpInstance, contactTracing, status]);

  if (status === 'SUCCESS') {
    return (
      <Redirect
        to={`/patient/${patientId}/pregnancy/${pregnancyId}/partners/contactTracing`}
      />
    );
  }

  if (patientData.loading) {
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
                `/patient/${patientId}/pregnancy/${pregnancyId}/partners/contactTracing`
              )
            }
          />
        </Box>
        <Box
          flex={false}
          direction={'row-responsive'}
          justify={'center'}
          align={'start'}
          fill={'horizontal'}
        >
          <Spinner />
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
              `/patient/${patientId}/pregnancy/${pregnancyId}/partners/contactTracing`
            )
          }
        />
        <Box
          flex={false}
          direction={'row-responsive'}
          justify={'center'}
          align={'start'}
          fill={'horizontal'}
        >
          <Box direction={'column'}>
            <Heading level={1} margin={'none'}>
              Contact Tracing
            </Heading>
            <Heading level={2} margin={'none'}>
              {`${patientData.data.patient.firstName} ${patientData.data.patient.lastName}`}
            </Heading>
          </Box>
        </Box>
        <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
          <Form onSubmit={onSubmit}>
            <FormField label={'Test'} name={'test'} required>
              <TextInput placeholder={'Test'} name={'test'} />
            </FormField>
            <FormField label={'Result'} name={'testResult'}>
              <TextInput placeholder={'Test Result'} name={'testResult'} />
            </FormField>
            <FormField label={'Comments'} name={'comments'}>
              <TextArea name={'comments'} placeholder={'Comments'} />
            </FormField>
            <FormField label={'Date'} name={'date'} required>
              <DateInput format={'yyyy-mm-dd'} name={'date'} />
            </FormField>
            <Box flex={false} align={'start'}>
              <Button type={'submit'} label={'Save'} primary />
            </Box>
          </Form>
        </Box>
      </Box>
    </Layout>
  );
};

export default CcontactTracingCreate;
