import {
  Box,
  DateInput,
  Form,
  FormField,
  TextInput,
  Text,
  Button,
} from 'grommet';
import React from 'react';
import { useHttpApi } from '../../../providers/HttpProvider';

const EditForm = ({ admission, closeForm }) => {
  const [facility, setFacility] = React.useState(admission.facility);
  const [dateAdmitted, setDateAdmitted] = React.useState(
    admission.dateAdmitted
  );
  const [admissionData, setAdmissionData] = React.useState();
  const { httpInstance } = useHttpApi();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');

  const onSubmit = (e) => {
    e.preventDefault();
    const data = {
      id: admission.id,
      patientId: admission.patientId,
      facility: facility,
      dateAdmitted: dateAdmitted,
    };
    setAdmissionData(data);
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const submit = async (body) => {
      try {
        await httpInstance.put(`/patient/hospitalAdmissions/${body.id}`, body);
        setStatus('SUCCESS');
      } catch (e) {
        console.error(e);
        setStatus('ERROR');
      }
    };

    if (status === 'SUBMIT' && admissionData) {
      submit(admissionData);
    }
    if (status === 'SUCCESS') {
      closeForm();
    }
  }, [httpInstance, status, admissionData, closeForm]);

  return (
    <Box>
      <Box flex={'grow'} overlfow={'auto'} pad={{ vertical: 'medium' }}>
        {status === 'ERROR' && (
          <Box
            fill={'horizontal'}
            pad={'medium'}
            gap={'medium'}
            background={'accent-2'}
          >
            <Text color={'dark-2'}>
              Oooooops! Error editing hospital admission!
            </Text>
          </Box>
        )}
        <Form onSubmit={onSubmit}>
          <FormField label={'Facility'} name={'facility'} required>
            <TextInput
              value={facility}
              name={'facility'}
              onChange={(e) => setFacility(e.target.value)}
            />
          </FormField>
          <FormField label={'Date Admitted'} name={'dateAdmitted'} required>
            <DateInput
              value={dateAdmitted}
              format={'yyyy-mm-dd'}
              name={'dateAdmitted'}
              onChange={(e) => setDateAdmitted(e.value)}
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
