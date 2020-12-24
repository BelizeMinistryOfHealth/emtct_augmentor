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

const PartnerSyphilisTreatmentCreate = () => {
  const [patientData, setPatientData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const [treatment, setTreatment] = React.useState();
  const { httpInstance } = useHttpApi();
  const { patientId, pregnancyId } = useParams();
  const history = useHistory();

  const onSubmit = (e) => {
    setTreatment({ ...e.value, patientId: parseInt(patientId) });
    setStatus('SUBMIT');
  };

  // Retrieve patient Information
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
          `/patients/${patientId}/pregnancy/${pregnancyId}/syphilisTreatments`,
          treatment
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
  }, [treatment, patientId, pregnancyId, status, httpInstance]);

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
              history.push(`/patient/${patientId}/partners/syphilisPartners`)
            }
          />
          <Box
            flex={false}
            direction={'row-responsive'}
            justify={'center'}
            align={'start'}
            fill={'horizontal'}
          >
            <Spinner />
          </Box>
        </Box>
      </Layout>
    );
  }

  if (status === 'SUCCESS') {
    return (
      <Redirect
        to={`/patient/${patientId}/pregnancy/${pregnancyId}/partners/syphilisTreatments`}
      />
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
              `/patient/${patientId}/pregnancy/${pregnancyId}/partners/syphilisTreatments`
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
              Create Syphilis Treatment for Partner
            </Heading>
            <Heading level={2} margin={'none'}>
              {`${patientData.data.patient.firstName} ${patientData.data.patient.lastName}`}
            </Heading>
          </Box>
        </Box>
        <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
          <Form onSubmit={onSubmit}>
            <FormField label={'Medication Name'} name={'medication'} required>
              <TextInput placeholder={'Medication Name'} name={'medication'} />
            </FormField>
            <FormField label={'Dosage'} name={'dosage'} required>
              <TextInput placeholder={'Dosage'} name={'dosage'} />
            </FormField>
            <FormField label={'Date'} name={'date'} required>
              <DateInput format={'yyyy-mm-dd'} name={'date'} />
            </FormField>
            <FormField label={'Comments'} name={'comments'}>
              <TextArea name={'comments'} placeholder={'Comments'} />
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

export default PartnerSyphilisTreatmentCreate;
