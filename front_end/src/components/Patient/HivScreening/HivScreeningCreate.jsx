import {
  Box,
  Heading,
  Button,
  Form,
  FormField,
  TextInput,
  DateInput,
} from 'grommet';
import { FormPreviousLink } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';

const HivScreeningCreateForm = () => {
  const [testName, setTestName] = React.useState();
  const [result, setResult] = React.useState();
  const [screeningDate, setScreningDate] = React.useState();
  const [dateSampleReceivedAtHq, setDateSampleReceivedAtHq] = React.useState();
  const [sampleCode, setSampleCode] = React.useState();
  const [destination, setDestination] = React.useState();
  const [dateResultReceived, setDateResultReceived] = React.useState();
  const [dateResultShared, setDateResultShared] = React.useState();
  const { patientId } = useParams();
  const history = useHistory();

  const onSubmit = () => {
    console.log('submitting');
  };

  return (
    <Box
      file={'vertical'}
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
      <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
        <Form onSubmit={onSubmit}>
          <FormField label={'Test Name'} name={'testName'} required>
            <TextInput
              value={testName}
              placeholder={'Test Name'}
              name={'testName'}
              onChange={(e) => setTestName(e.target.value)}
            />
          </FormField>
          <FormField label={'Result'} name={'result'} required>
            <TextInput
              value={result}
              placeholder={'Test Result'}
              name={'result'}
              onChange={(e) => setResult(e.target.value)}
            />
          </FormField>
          <FormField label={'Screening Date'} name={'screeningDate'} required>
            <DateInput
              format={'yyyy-mm-dd'}
              name={'screeningDate'}
              value={screeningDate}
            />
          </FormField>
          <FormField
            label={'Date Sample Received At HQ'}
            name={'dateSampleReceivedAtHq'}
            required
          >
            <DateInput
              format={'yyyy-mm-dd'}
              name={'dateSampleReceivedAtHq'}
              value={dateSampleReceivedAtHq}
            />
          </FormField>
          <FormField
            label={'Date Result Recieved'}
            name={'dateResultReceived'}
            required
          >
            <DateInput
              format={'yyyy-mm-dd'}
              name={'dateResultReceived'}
              value={dateResultReceived}
            />
          </FormField>
          <FormField
            label={'Date Result Shared'}
            name={'dateResultShared'}
            required
          >
            <DateInput
              format={'yyyy-mm-dd'}
              name={'dateResultShared'}
              value={dateResultShared}
            />
          </FormField>
          <FormField label={'Sample Code'} name={'sampleCode'} required>
            <TextInput
              value={sampleCode}
              placeholder={'Sample Code'}
              name={'sampleCode'}
              onChange={(e) => setSampleCode(e.target.value)}
            />
          </FormField>
          <FormField label={'Destination'} name={'destination'} required>
            <TextInput
              value={destination}
              placeholder={'Destination'}
              name={'destination'}
              onChange={(e) => setDestination(e.target.value)}
            />
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
