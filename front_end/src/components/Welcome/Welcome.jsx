import React from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import { Box, Heading, Button } from 'grommet';

const LoginButton = () => {
  const { isAuthenticated, loginWithRedirect } = useAuth0();

  return (
    !isAuthenticated && <Button label='Sign In' onClick={loginWithRedirect} />
  );
};

const Welcome = () => {
  return (
    <Box justify={'center'} align={'center'} background={'brand'} fill>
      <Heading>Welcome to MCH EMTCT Project</Heading>
      <Box align='center' pad='medium'>
        <LoginButton />
      </Box>
    </Box>
  );
};

export default Welcome;
