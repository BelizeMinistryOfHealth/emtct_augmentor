import {
  Box,
  Heading,
  Button,
  Form,
  FormField,
  TextInput,
  DateInput,
  Text,
} from 'grommet';
import { FormPreviousLink } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import Layout from '../../Layout/Layout';

const HivScreeningCreateForm = () => {
  const [screening, setScreening] = React.useState();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const [patientData, setPatientData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });
  const [errorMessage, setErrorMessage] = React.useState(undefined);
  const { patientId } = useParams();
  const history = useHistory();
  const { httpInstance } = useHttpApi();

  const onSubmit = (e) => {
    setScreening({ ...e.value, patientId: parseInt(patientId) });
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const fetchPatient = async () => {
      try {
        const result = await httpInstance.get(`/patient/${patientId}`);
        console.log({ result: result.data });
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
    const post = (screening) => {
      const mchEncounterId = patientData.data.antenatalEncounter.id;
      httpInstance
        .post(`/patient/hivScreening`, {
          ...screening,
          mchEncounterId,
        })
        .then(() => {
          setStatus('SUCCESS');
          setErrorMessage(undefined);
        })
        .catch((e) => {
          console.error(e);
          setStatus('ERROR');
          if (e.response) {
            setErrorMessage(e.response.data);
          }
        });
    };
    if (status === 'SUBMIT') {
      setErrorMessage(undefined);
      post(screening);
    }
  }, [screening, httpInstance, status, patientData, errorMessage]);

  if (status === 'SUBMIT') {
    return (
      <Box
        fill={'vertical'}
        overflow={'auto'}
        pad={'medium'}
        width={'xlarge'}
        justify={'center'}
      >
        <Text size={'xlarge'}>Saving...</Text>
      </Box>
    );
  }

  if (status === 'SUCCESS') {
    history.push(`/patient/${patientId}/hiv_screenings`);
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
          onClick={() => history.push(`/patient/${patientId}/hiv_screenings`)}
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
              Create HIV Screening
            </Heading>
            {patientData.data && (
              <Heading level={3} margin={'none'}>
                {`${patientData.data.patient.firstName} ${patientData.data.patient.lastName}`}
              </Heading>
            )}
          </Box>
        </Box>
        {status === 'ERROR' && (
          <Box
            fill={'horizontal'}
            pad={'medium'}
            gap={'medium'}
            background={'accent-4'}
          >
            {errorMessage ? (
              <Text>{errorMessage}</Text>
            ) : (
              <Text>Error creating hiv screening!</Text>
            )}
          </Box>
        )}

        <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
          <Form onSubmit={onSubmit}>
            <FormField label={'Test Name'} name={'testName'} required>
              <TextInput placeholder={'Test Name'} name={'testName'} />
            </FormField>
            <FormField label={'Result'} name={'result'}>
              <TextInput placeholder={'Test Result'} name={'result'} />
            </FormField>
            <FormField label={'Screening Date'} name={'screeningDate'} required>
              <DateInput format={'yyyy-mm-dd'} name={'screeningDate'} />
            </FormField>
            <FormField
              label={'Date Sample Received At HQ'}
              name={'dateSampleReceivedAtHq'}
            >
              <DateInput
                format={'yyyy-mm-dd'}
                name={'dateSampleReceivedAtHq'}
              />
            </FormField>
            <FormField
              label={'Date Sample Taken'}
              name={'dateSampleTaken'}
              required
            >
              <DateInput format={'yyyy-mm-dd'} name={'dateSampleTaken'} />
            </FormField>
            <FormField
              label={'Date Result Received'}
              name={'dateResultReceived'}
            >
              <DateInput format={'yyyy-mm-dd'} name={'dateResultReceived'} />
            </FormField>
            <FormField label={'Date Result Shared'} name={'dateResultShared'}>
              <DateInput format={'yyyy-mm-dd'} name={'dateResultShared'} />
            </FormField>
            <FormField label={'Date Sample Shipped'} name={'dateSampleShipped'}>
              <DateInput format={'yyyy-mm-dd'} name={'dateSampleShipped'} />
            </FormField>
            <FormField label={'Sample Code'} name={'sampleCode'} required>
              <TextInput placeholder={'Sample Code'} name={'sampleCode'} />
            </FormField>
            <FormField label={'Destination'} name={'destination'} required>
              <TextInput placeholder={'Destination'} name={'destination'} />
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

export default HivScreeningCreateForm;
