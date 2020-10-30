import {
  Box,
  Button,
  DateInput,
  Form,
  FormField,
  Heading,
  Text,
  TextInput,
} from 'grommet';
import { FormPreviousLink } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';

const HospitalAdmissionCreateForm = () => {
  const [admission, setAdmission] = React.useState();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const { patientId } = useParams();
  const history = useHistory();
  const { httpInstance } = useHttpApi();

  const onSubmit = (e) => {
    setAdmission({ ...e.value, patientId: parseInt(patientId) });
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const post = async (admission) => {
      try {
        await httpInstance.post(`/patient/hospitalAdmissions`, admission);
        setStatus('SUBMIT');
      } catch (e) {
        console.error(e);
        setStatus('ERROR');
      }
    };
    if (status === 'SUBMIT') {
      post(admission);
    }
  }, [httpInstance, admission, status]);

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
          onClick={() => history.push(`/patient/${patientId}/admissions`)}
        />
        <Box
          flex={false}
          direction={'row-responsive'}
          justify={'center'}
          fill={'horizontal'}
        >
          <Heading level={2} margin={'none'}>
            Successfully Saved Hospital Admission!
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
        onClick={() => history.push(`/patient/${patientId}/admissions`)}
      />
      <Box
        flex={false}
        direction={'row-responsive'}
        justify={'center'}
        fill={'horizontal'}
      >
        <Heading level={2} margin={'none'}>
          Create Hospital Admission
        </Heading>
      </Box>
      {status === 'ERROR' && (
        <Box
          fill={'horizontal'}
          pad={'medium'}
          gap={'medium'}
          background={'red'}
        >
          <Text>Error creating hospital admission!</Text>
        </Box>
      )}

      <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
        <Form onSubmit={onSubmit}>
          <FormField label={'Facility'} name={'facility'} required>
            <TextInput placeholder={'Facility'} name={'facility'} />
          </FormField>
          <FormField label={'Date Admitted'} name={'dateAdmitted'} required>
            <DateInput format={'yyyy-mm-dd'} name={'dateAdmitted'} />
          </FormField>
          <Box flex={false} align={'center'}>
            <Button type={'submit'} label={'Save'} primary />
          </Box>
        </Form>
      </Box>
    </Box>
  );
};

export default HospitalAdmissionCreateForm;
