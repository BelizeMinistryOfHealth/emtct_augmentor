import {
  Box,
  Button,
  DateInput,
  Form,
  FormField,
  Text,
  TextArea,
  TextInput,
} from 'grommet';
import React from 'react';
import { useHttpApi } from '../../../providers/HttpProvider';

const ContactTracingEdit = ({ contactTracing, closeEditScreen }) => {
  const [test, setTest] = React.useState(contactTracing.test);
  const [testResult, setTestResult] = React.useState(contactTracing.testResult);
  const [comments, setComments] = React.useState(contactTracing.comments);
  const [date, setDate] = React.useState(contactTracing.date);
  const [contactTracingData, setContactTracingData] = React.useState();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const { httpInstance } = useHttpApi();

  const onSubmit = (e) => {
    e.preventDefault();
    const data = {
      ...contactTracing,
      test,
      testResult,
      comments,
      date,
    };
    setContactTracingData(data);
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const submit = () => {
      httpInstance
        .put(
          `/patient/${contactTracingData.patientId}/partners/contactTracing`,
          contactTracingData
        )
        .then(() => {
          setStatus('SUCCESS');
        })
        .catch((e) => {
          console.error(e);
          setStatus('ERROR');
        });
    };
    if (status === 'SUBMIT' && contactTracingData) {
      submit();
    }
  }, [contactTracingData, httpInstance, status]);

  if (status === 'SUCCESS') {
    closeEditScreen();
  }

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
            <Text color={'dark-2'}>
              Ooooops! Error editing contact tracing!
            </Text>
          </Box>
        )}
        <Form onSubmit={onSubmit}>
          <FormField labe={'Test'} name={'test'} required>
            <TextInput
              value={test}
              name={'test'}
              onSubmit={(e) => setTest(e.target.value.trim())}
              onChange={(e) => setTest(e.target.value)}
            />
          </FormField>
          <FormField label={'Test Result'} name={'testResult'}>
            <TextInput
              value={testResult}
              name={'testResult'}
              onSubmit={(e) => setTestResult(e.target.value.trim())}
              onChange={(e) => setTestResult(e.target.value)}
            />
          </FormField>
          <FormField label={'Comments'} name={'comments'}>
            <TextArea
              value={comments}
              name={'comments'}
              onSubmit={(e) => setComments(e.target.value.trim())}
              onChange={(e) => setComments(e.target.value)}
            />
          </FormField>
          <FormField label={'Date'} name={'date'}>
            <DateInput
              format={'yyyy-mm-dd'}
              calendarProps={{ locale: 'en-BZ', fill: false }}
              value={date}
              name={'date'}
              onChange={(e) => {
                setDate(e.value);
              }}
            />
          </FormField>
          <Box flex={false} align={'start'}>
            <Button type={'submit'} label={'Save'} primary />
          </Box>
        </Form>
      </Box>
    </Box>
  );
};

export default ContactTracingEdit;
