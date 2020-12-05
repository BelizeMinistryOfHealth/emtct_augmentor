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

const EditForm = ({ contraceptive, onCloseForm }) => {
  const [name, setName] = React.useState(contraceptive.contraceptive);
  const [comments, setComments] = React.useState(contraceptive.comments);
  const [dateUsed, setDateUsed] = React.useState(contraceptive.dateUsed);
  const [contraceptiveData, setContraceptiveData] = React.useState();
  const { httpInstance } = useHttpApi();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');

  const onSubmit = (e) => {
    e.preventDefault();
    const data = {
      patientId: contraceptive.patientId,
      id: contraceptive.id,
      contraceptive: name,
      comments,
      dateUsed,
    };
    setContraceptiveData(data);
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const submit = async (body) => {
      try {
        await httpInstance.put(`/contraceptivesUsed`, body);
        setStatus('SUCCESS');
        onCloseForm();
      } catch (e) {
        console.error(e);
        setStatus('ERROR');
      }
    };
    if (status === 'SUBMIT' && contraceptiveData) {
      submit(contraceptiveData);
    }
  }, [contraceptiveData, status, httpInstance, onCloseForm]);

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
            <Text color={'dark-2'}>Ooooops! Error editing contraceptive!</Text>
          </Box>
        )}
        <Form onSubmit={onSubmit}>
          <FormField label={'Contraceptive'} name={'name'} required>
            <TextInput
              value={name}
              name={'name'}
              onChange={(e) => setName(e.target.value)}
            />
          </FormField>
          <FormField label={'Comments'} name={'comments'}>
            <TextArea
              value={comments}
              name={'comments'}
              onChange={(e) => setComments(e.target.value)}
            />
          </FormField>
          <FormField label={'Date used'} name={'dateUsed'} required>
            <DateInput
              format={'yyyy-mm-dd'}
              value={dateUsed}
              name={'dateUsed'}
              onChange={(e) => setDateUsed(e.value)}
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
