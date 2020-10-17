import React from 'react';
import { Box, Button, Card, Form, FormField, TextInput } from 'grommet';

const PatientIdSearch = (props) => {
  const [patientId, setPatientId] = React.useState();
  const { onSubmit } = props;

  return (
    <Card
      pad={'small'}
      gap={'medium'}
      background={'light-1'}
      margin={{ right: '20' }}
    >
      <Box width={'medium'}>
        <Form onSubmit={() => onSubmit(patientId)}>
          <FormField
            label={'Patient Id'}
            name={'patientId'}
            validate={{ regexp: /^[0-9]{4,20}$/, message: '4-20 digits' }}
            required
          >
            <TextInput
              name={'patientId'}
              type={'name'}
              onChange={(change) => {
                const pId = parseInt(change.target.value);
                if (!isNaN(pId)) {
                  setPatientId(pId);
                }
              }}
            />
          </FormField>
          <Box direction={'row'} justify={'between'} margin={{ top: 'medium' }}>
            <Button type={'submit'} label={'Search'} />
          </Box>
        </Form>
      </Box>
    </Card>
  );
};

export default PatientIdSearch;
