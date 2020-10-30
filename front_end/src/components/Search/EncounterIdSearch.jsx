import { Box, Button, Card, Form, FormField, TextInput } from 'grommet';
import React from 'react';

//TODO:  Add onSubmit Handler
const EncounterIdSearch = () => {
  const [, setEncounterNumber] = React.useState();
  return (
    <Card
      pad={'small'}
      gap={'medium'}
      background={'light-1'}
      margin={{ left: '20' }}
    >
      <Box width={'medium'} margin={{ left: 'medium' }}>
        <Form>
          <FormField
            label={'Encounter #'}
            name={'encounterNumber'}
            validate={{ regexp: /^[0-9]{4,20}$/, message: '4-20 digits' }}
            required
          >
            <TextInput
              name={'encounterNumber'}
              type={'name'}
              onChange={(e) => {
                const eId = parseInt(e.target.value);
                if (!isNaN(eId)) {
                  setEncounterNumber(eId);
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

export default EncounterIdSearch;
