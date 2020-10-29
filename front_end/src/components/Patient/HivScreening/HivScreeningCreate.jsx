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

const HivScreeningCreateForm = () => {
  const [screening, setScreening] = React.useState();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const { patientId } = useParams();
  const history = useHistory();
  const { httpInstance } = useHttpApi();

  const onSubmit = (e) => {
    setScreening({ ...e.value, patientId: parseInt(patientId) });
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const post = async (screening) => {
      try {
        await httpInstance.post(`/patient/hivScreening`, screening);
        setStatus('SUCCESS');
      } catch (e) {
        console.error(e);
        setStatus('ERROR');
      }
    };
    if (status === 'SUBMIT') {
      post(screening);
    }
  }, [screening, httpInstance, status]);

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
        fill={'horizontal'}
      >
        <Heading level={2} margin={'non'}>
          Create HIV Screening
        </Heading>
      </Box>
      {status === 'ERROR' && (
        <Box
          fill={'horizontal'}
          pad={'medium'}
          gap={'medium'}
          background={'red'}
        >
          <Text>Error creating hiv screening!</Text>
        </Box>
      )}

      <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
        <Form onSubmit={onSubmit}>
          <FormField label={'Test Name'} name={'testName'} required>
            <TextInput placeholder={'Test Name'} name={'testName'} />
          </FormField>
          <FormField label={'Result'} name={'result'} required>
            <TextInput placeholder={'Test Result'} name={'result'} />
          </FormField>
          <FormField label={'Screening Date'} name={'screeningDate'} required>
            <DateInput format={'yyyy-mm-dd'} name={'screeningDate'} />
          </FormField>
          <FormField
            label={'Date Sample Received At HQ'}
            name={'dateSampleReceivedAtHq'}
            required
          >
            <DateInput format={'yyyy-mm-dd'} name={'dateSampleReceivedAtHq'} />
          </FormField>
          <FormField label={'Date Result Received'} name={'dateResultReceived'}>
            <DateInput format={'yyyy-mm-dd'} name={'dateResultReceived'} />
          </FormField>
          <FormField label={'Date Result Shared'} name={'dateResultShared'}>
            <DateInput format={'yyyy-mm-dd'} name={'dateResultShared'} />
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
  );
};

export default HivScreeningCreateForm;
