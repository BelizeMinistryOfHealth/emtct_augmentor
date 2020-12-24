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

const PartnerSyphilisTreatmentEdit = ({ treatment, closeEditScreen }) => {
  const [medication, setMedication] = React.useState(treatment.medication);
  const [dosage, setDosage] = React.useState(treatment.dosage);
  const [comments, setComments] = React.useState(treatment.comments);
  const [date, setDate] = React.useState(treatment.date);
  const [treatmentData, setTreatmentData] = React.useState();
  const { httpInstance } = useHttpApi();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');

  const onSubmit = (e) => {
    e.preventDefault();
    const data = {
      ...treatment,
      medication,
      dosage,
      comments,
      date,
    };
    setTreatmentData(data);
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const submit = () => {
      httpInstance
        .put(
          `/patients/${treatmentData.patientId}/pregnancy/${treatmentData.pregnancyId}/syphilisTreatments`,
          treatmentData
        )
        .then(() => {
          setStatus('SUCCESS');
        })
        .catch((e) => {
          console.error(e);
          setStatus('ERROR');
        });
    };
    if (status === 'SUBMIT' && treatmentData) {
      submit();
    }
  }, [treatmentData, httpInstance, status]);

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
              Oooops! Error editing syphilis treatment!
            </Text>
          </Box>
        )}
        <Form onSubmit={onSubmit}>
          <FormField label={'Medication'} name={'medication'} required>
            <TextInput
              value={medication}
              name={'medication'}
              onChange={(e) => setMedication(e.target.value.trim())}
            />
          </FormField>
          <FormField label={'Dosage'} name={'dosage'} required>
            <TextInput
              value={dosage}
              name={'dosage'}
              onChange={(e) => setDosage(e.target.value.trim())}
            />
          </FormField>
          <FormField label={'Date'} name={'date'} required>
            <DateInput
              format={'yyyy-mm-dd'}
              value={date}
              name={'date'}
              onChange={(e) => setDate(e.value)}
            />
          </FormField>
          <FormField label={'Comments'} name={'comments'}>
            <TextArea
              name={'comments'}
              value={comments}
              onChange={(e) => setComments(e.target.value)}
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

export default PartnerSyphilisTreatmentEdit;
