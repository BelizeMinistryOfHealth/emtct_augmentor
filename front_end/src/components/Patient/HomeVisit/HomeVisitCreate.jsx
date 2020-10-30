import {
  Box,
  FormField,
  TextArea,
  Heading,
  Button,
  DateInput,
  Form,
} from 'grommet';
import { FormPreviousLink } from 'grommet-icons';
import React from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';

const HomeVisitCreateForm = () => {
  const [reason, setReason] = React.useState();
  const [comments, setComments] = React.useState();
  const [dateOfVisit, setDateOfVisit] = React.useState();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const [, setHomeVisit] = React.useState();
  const { patientId } = useParams();
  const history = useHistory();

  const { httpInstance } = useHttpApi();

  const onChangeReason = (value) => setReason(value);
  const onChangeComments = (value) => setComments(value);
  const onChangeDateOfVisit = setDateOfVisit;

  const onSubmit = (e) => {
    e.preventDefault();
    setStatus('SUBMIT');
  };

  React.useEffect(() => {
    const submit = async () => {
      const body = {
        reason,
        comments,
        patientId: parseInt(patientId),
        dateOfVisit,
      };

      try {
        const result = await httpInstance.post('/patient/homeVisits', body);
        setHomeVisit(result.data);
        setStatus('SUCCESS');
      } catch (e) {
        console.error(e);
        setStatus('ERROR');
      }
    };
    if (status === 'SUBMIT') {
      submit();
    }
  }, [status, httpInstance, reason, comments, patientId, dateOfVisit]);

  if (status === 'SUCCESS') {
    history.push(`/patient/${patientId}/home_visits`);
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
        onClick={() => history.push(`/patient/${patientId}/home_visits`)}
      />
      <Box
        flex={false}
        direction={'row-responsive'}
        justify={'center'}
        fill={'horizontal'}
      >
        <Heading level={2} margin={'none'}>
          Create Home Visit
        </Heading>
      </Box>
      <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
        <Form onSubmit={onSubmit}>
          <FormField label={'Reason'} name={'reason'} required>
            <TextArea
              value={reason}
              placeholder={'Reason'}
              name={'reason'}
              onChange={(e) => onChangeReason(e.target.value)}
            />
          </FormField>
          <FormField label={'Comments'} name={'comments'} required>
            <TextArea
              value={comments}
              name={'comments'}
              placeholder={'Comments'}
              onChange={(e) => onChangeComments(e.target.value)}
            />
          </FormField>
          <FormField label={'Date of Visit'} name={'dateOfVisit'} required>
            <DateInput
              format={'yyyy-mm-dd'}
              name={'dateOfVisit'}
              value={dateOfVisit}
              onChange={(e) => {
                console.dir(e);
                onChangeDateOfVisit(e.value);
              }}
            />
          </FormField>
          <Box flex={false} as={'footer'} align={'start'}>
            <Button type={'submit'} label={'Save'} primary />
          </Box>
        </Form>
      </Box>
    </Box>
  );
};

export default HomeVisitCreateForm;
