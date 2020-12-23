import {
  Box,
  FormField,
  TextArea,
  Text,
  Heading,
  Button,
  DateInput,
  Form,
} from 'grommet';
import { FormPreviousLink } from 'grommet-icons';
import React from 'react';
import { useParams, useHistory } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import Layout from '../../Layout/Layout';

const HomeVisitCreateForm = () => {
  const [reason, setReason] = React.useState();
  const [comments, setComments] = React.useState();
  const [dateOfVisit, setDateOfVisit] = React.useState();
  // Form status: START -> SUBMIT -> ERROR -> SUCCESS
  const [status, setStatus] = React.useState('START');
  const [, setHomeVisit] = React.useState();
  const { patientId, pregnancyId } = useParams();
  const [patientData, setPatientData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });
  const [errorMessage, setErrorMessage] = React.useState(undefined);
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
    const fetchPatient = async () => {
      try {
        const result = await httpInstance.get(`/patients/${patientId}`);
        setPatientData({ data: result.data, loading: false, error: undefined });
      } catch (e) {
        console.error(e);
        setPatientData({ data: undefined, loading: false, error: e });
      }
    };
    if (patientData.loading) {
      fetchPatient();
    }
  }, [httpInstance, patientId, patientData]);

  React.useEffect(() => {
    const submit = () => {
      const body = {
        reason,
        comments,
        patientId: parseInt(patientId),
        dateOfVisit,
      };

      httpInstance
        .post(
          `/patients/${patientId}/pregnancy/${pregnancyId}/homeVisits`,
          body
        )
        .then((response) => {
          setHomeVisit(response.data);
          setStatus('SUCCESS');
        })
        .catch((e) => {
          console.error(e);
          setStatus('ERROR');
          if (e.response) {
            setErrorMessage(e.response.data);
          }
        });
    };
    if (status === 'SUBMIT') {
      submit();
    }
  }, [
    status,
    httpInstance,
    reason,
    comments,
    patientId,
    pregnancyId,
    dateOfVisit,
    errorMessage,
    patientData,
  ]);

  if (status === 'SUCCESS') {
    history.push(`/patient/${patientId}/pregnancy/${pregnancyId}/home_visits`);
  }

  return (
    <Layout>
      <Box
        fill={'vertical'}
        overflow={'auto'}
        pad={'medium'}
        width={'xlarge'}
        justify={'center'}
      >
        <Button
          icon={<FormPreviousLink size={'large'} />}
          onClick={() =>
            history.push(
              `/patient/${patientId}/pregnancy/${pregnancyId}/home_visits`
            )
          }
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
        {status === 'ERROR' && (
          <Box
            fill={'horizontal'}
            pad={'medium'}
            gap={'medium'}
            background={'accent-4'}
          >
            {errorMessage ? (
              <Text>{errorMessage}</Text>
            ) : (
              <Text>Error creating a new home visit!</Text>
            )}
          </Box>
        )}

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
    </Layout>
  );
};

export default HomeVisitCreateForm;
