import React from 'react';
import { Box, CardHeader, Text } from 'grommet';

const AppCardHeader = ({ title, patient, ...rest }) => {
  return (
    <CardHeader {...rest}>
      <Box direction={'row'} align={'start'} fill='horizontal'>
        <Box
          direction={'column'}
          align={'start'}
          fill={'horizontal'}
          justify={'between'}
          alignContent={'center'}
        >
          <Text size={'xxlarge'} weight={'bold'} textAlign={'start'}>
            {title}
          </Text>
          <Text size={'large'} textAlign={'end'} weight={'normal'}>
            {patient.firstName} {patient.lastName}{' '}
          </Text>
        </Box>
      </Box>
    </CardHeader>
  );
};

export default AppCardHeader;
