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

const ContraceptivesCreateForm = () => {
  const [contraceptive, setContraceptive] = React.useState();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const { patientId } = useParams();
  const history = useHistory();
  const { httpInstance } = useHttpApi();

  const onSubmit = (e) => {
    setContraceptive({ ...e.value, patientId: parseInt(patientId) });
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const post = async (contraceptive) => {
      try {
        await httpInstance.post(`/patient/contraceptivesUsed`, contraceptive);
        setStatus('SUCCESS');
      } catch (e) {
        console.error(e);
        setStatus('ERROR');
      }
    };
    if (status === 'SUBMIT') {
      post(contraceptive);
    }
  }, [contraceptive, httpInstance, status]);

  if (status === 'SUCCESS') {
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
          onClick={() => history.push(`/patient/${patientId}/contraceptives`)}
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
    );
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
        onClick={() => history.push(`/patient/${patientId}/contraceptives`)}
      ></Button>
      <Box
        flex={false}
        direction={'row-responsive'}
        justify={'center'}
        fill={'horizontal'}
      >
        <Heading level={2} margin={'none'}>
          Create Contraceptive Usage
        </Heading>
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
  );
};

export default ContraceptivesCreateForm;
