import React from 'react';
import { Box, Button, Heading, Text } from 'grommet';

const UserDelete = ({ onClickNo, onClickYes }) => {
  return (
    <Box pad={'medium'} gap={'small'} width={'medium'}>
      <Heading level={3} margin={'none'}>
        Confirm
      </Heading>
      <Text>Are you sure you want to delete this user?</Text>
      <Box
        as='footer'
        gap='small'
        direction='row'
        align='center'
        justify='end'
        pad={{ top: 'medium', bottom: 'small' }}
      >
        <Button label={'No'} color={'dark-3'} onClick={onClickNo} />
        <Button
          label={
            <Text color='white'>
              <strong>Yes</strong>
            </Text>
          }
          primary
          color={'status-critical'}
          onClick={onClickYes}
        />
      </Box>
    </Box>
  );
};

export default UserDelete;
