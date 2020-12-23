import { Box, Button, DateInput, FormField, TextArea, Text } from 'grommet';
import React from 'react';
import { useHttpApi } from '../../../providers/HttpProvider';
import { useParams } from 'react-router-dom';

const EditForm = ({ visit, closeForm }) => {
  const [reason, setReason] = React.useState(visit.reason);
  const [comments, setComments] = React.useState(visit.comments);
  const [dateOfVisit, setDateOfVisit] = React.useState(visit.dateOfVisit);
  const [submitting, setSubmitting] = React.useState(false);
  const { httpInstance } = useHttpApi();
  const [error, setError] = React.useState(undefined);
  const [homeVisit, setHomeVisit] = React.useState(undefined);
  const [editStatus, setEditStatus] = React.useState(undefined);

  const { patientId, pregnancyId } = useParams();

  const onChangeReason = (value) => {
    setReason(value);
  };
  const onChangeComments = (value) => {
    setComments(value);
  };
  const onChangeDateOfVisit = (value) => {
    setDateOfVisit(value);
  };

  React.useEffect(() => {
    const edit = async () => {
      try {
        await httpInstance.put(
          `/patients/${patientId}/pregnancy/${pregnancyId}/homeVisits`,
          homeVisit
        );
        setSubmitting(false);
        setEditStatus('Successfully edited');
        closeForm();
      } catch (e) {
        console.error(e);
        setError('Edit request failed!');
        setSubmitting(false);
        setEditStatus('Error while editing');
      }
    };
    if (homeVisit && submitting) {
      edit();
    }
  }, [
    httpInstance,
    submitting,
    homeVisit,
    closeForm,
    editStatus,
    patientId,
    pregnancyId,
  ]);

  const onSubmit = (e) => {
    e.preventDefault();
    const data = {
      ...visit,
      reason,
      comments,
      dateOfVisit,
    };
    setSubmitting(true);
    setHomeVisit(data);
  };

  if (submitting) {
    return <>Submitting....</>;
  }

  return (
    <>
      <Box flex={'grow'} overflow={'auto'} pad={{ vertical: 'medium' }}>
        {error && <>${editStatus}</>}
        {!error && editStatus && (
          <>
            <Text>{editStatus}</Text>
          </>
        )}
        <FormField label={'Reason'} required>
          <TextArea
            value={reason}
            placeholder={'Reason'}
            onChange={(e) => onChangeReason(e.target.value)}
          />
        </FormField>
        <FormField label={'Comments'} required>
          <TextArea
            value={comments}
            placeholder={'Comments'}
            onChange={(e) => onChangeComments(e.target.value)}
          />
        </FormField>
        <FormField label={'Date of Visit'} required>
          <DateInput
            format={'yyyy-mm-dd'}
            value={dateOfVisit}
            onChange={(e) => {
              console.dir(e);
              onChangeDateOfVisit(e.value);
            }}
          />
        </FormField>
      </Box>
      <Box flex={false} as={'footer'} align={'start'}>
        <Button type={'submit'} label={'Edit'} onClick={onSubmit} primary />
      </Box>
    </>
  );
};

export default EditForm;
