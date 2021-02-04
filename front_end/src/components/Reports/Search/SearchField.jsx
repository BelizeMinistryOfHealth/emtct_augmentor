import React from 'react';
import { Box, Button, Card, Form, FormField, TextInput } from 'grommet';

const SearchField = ({ onSubmit, label }) => {
  const [query, setQuery] = React.useState();
  return (
    <Card
      pad={'small'}
      gap={'medium'}
      background={'light-1'}
      margin={{ right: '20' }}
    >
      <Box width={'medium'}>
        <Form onSubmit={() => onSubmit(query)}>
          <FormField label={label} name={'query'} required>
            <TextInput
              name={'query'}
              onChange={(e) => setQuery(e.target.value.trim())}
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

export default SearchField;
