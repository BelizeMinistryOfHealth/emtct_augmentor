import React from 'react';
import { Anchor, Box, Header } from 'grommet';
import { withAuth0 } from '@auth0/auth0-react';
import Logout from '../Auth/Logout';

const Navbar = () => {
  return (
    <Header background={'dark-1'} pad={'small'}>
      <Box direction={'row'} align={'center'} gap={'small'}>
        <Anchor color={'white'} href={'/'}>
          EMTCT Project
        </Anchor>
      </Box>
      <Box direction={'row'} align={'end'} gap={'small'}>
        <Logout />
      </Box>
    </Header>
  );
};

export default withAuth0(Navbar);
