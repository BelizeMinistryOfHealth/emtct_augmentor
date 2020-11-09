import {
  FormField,
  Form,
  Box,
  TextInput,
  DateInput,
  Button,
  Text,
} from 'grommet';
import React from 'react';
import { useHttpApi } from '../../../providers/HttpProvider';

const EditForm = ({ screening }) => {
  const [testName, setTestName] = React.useState(screening.testName);
  const [result, setResult] = React.useState(screening.result);
  const [sampleCode, setSampleCode] = React.useState(screening.sampleCode);
  const [destination, setDestination] = React.useState(screening.destination);
  const [screeningDate, setScreeningDate] = React.useState(
    screening.screeningDate
  );
  const [dateSampleReceivedAtHq, setDateSampleReceivedAtHq] = React.useState(
    screening.dateSampleReceivedAtHq
  );
  const [dateResultReceived, setDateResultReceived] = React.useState(
    screening.dateResultReceived
  );
  const [dateResultShared, setDateResultShared] = React.useState(
    screening.dateResultShared
  );
  const [dateSampleTaken, setDateSampleTaken] = React.useState(
    screening.dateSampleTaken
  );

  const [screenData, setScreenData] = React.useState();
  const { httpInstance } = useHttpApi();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');

  const onSubmit = (e) => {
    e.preventDefault();
    const data = {
      patientId: screening.patientId,
      id: screening.id,
      result,
      testName,
      destination,
      sampleCode,
      screeningDate,
      dateSampleReceivedAtHq,
      dateResultReceived,
      dateResultShared,
      dateSampleTaken,
    };
    setScreenData(data);
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const submit = async (body) => {
      try {
        await httpInstance.put(`/patient/hivScreening/${body.id}`, body);
        setStatus('SUCCESS');
      } catch (e) {
        console.error(e);
        setStatus('ERROR');
      }
    };
    if (status === 'SUBMIT' && screenData) {
      submit(screenData);
    }
  }, [screenData, status, httpInstance]);

  return (
    <Box>
      <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
        {status === 'ERROR' && (
          <Box
            fill={'horizontal'}
            pad={'medium'}
            gap={'medium'}
            background={'accent-2'}
          >
            <Text color={'dark-2'}>Ooops! Error editing hiv sreening!</Text>
          </Box>
        )}
        <Form onSubmit={onSubmit}>
          <FormField label={'Test Name'} name={'testName'} required>
            <TextInput
              value={testName}
              name={'testName'}
              onChange={(e) => setTestName(e.target.value.trim())}
            />
          </FormField>
          <FormField label={'Result'} name={'result'}>
            <TextInput
              value={result}
              name={'result'}
              onChange={(e) => setResult(e.target.value.trim())}
            />
          </FormField>
          <FormField label={'Sample Code'} name={'sampleCode'} required>
            <TextInput
              value={sampleCode}
              name={'sampleCode'}
              onChange={(e) => setSampleCode(e.target.value.trim())}
            />
          </FormField>
          <FormField label={'Destination'} name={'destination'} required>
            <TextInput
              value={destination}
              name={'destination'}
              onChange={(e) => setDestination(e.target.value.trim())}
            />
          </FormField>
          <FormField label={'Screening Date'} name={'screeningDate'} required>
            <DateInput
              format={'yyyy-mm-dd'}
              value={screeningDate}
              name={'screeningDate'}
              onChange={(e) => setScreeningDate(e.value)}
            />
          </FormField>
          <FormField
            label={'Date Sample Received at HQ'}
            name={'dateSampleReceivedAtHq'}
          >
            <DateInput
              format={'yyyy-mm-dd'}
              value={dateSampleReceivedAtHq}
              name={'dateSampleReceivedAtHq'}
              onChange={(e) => setDateSampleReceivedAtHq(e.value)}
            />
          </FormField>
          <FormField
            label={'Date Sample Taken'}
            name={'dateSampleTaken'}
            required
          >
            <DateInput
              format={'yyyy-mm-dd'}
              name={'dateSampleTaken'}
              value={dateSampleTaken}
              onChange={(e) => setDateSampleTaken(e.value)}
            />
          </FormField>
          <FormField label={'Date Result Receivd'} name={'dateResultReceived'}>
            <DateInput
              format={'yyyy-mm-dd'}
              value={dateResultReceived}
              name={'dateResultReceived'}
              onChange={(e) => setDateResultReceived(e.value)}
            />
          </FormField>
          <FormField label={'Date Result Shared'} name={'dateResultShared'}>
            <DateInput
              format={'yyyy-mm-dd'}
              value={dateResultShared}
              name={'dateResultShared'}
              onChange={(e) => setDateResultShared(e.value)}
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

export default EditForm;
