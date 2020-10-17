import React from 'react';
import { Anchor, Box, Header } from 'grommet';
import { withAuth0 } from '@auth0/auth0-react';
import LogoutButton from '../Logout/Logout';

const Navbar = (props) => {
  const { user, isAuthenticated } = props.auth0;

  return (
    isAuthenticated && (
      <Header background={'dark-1'} pad={'small'}>
        <Box direction={'row'} align={'center'} gap={'small'}>
          <Anchor color={'white'}>EMTCT Project | {user.name}</Anchor>
        </Box>
        <Box direction={'row'} align={'end'} gap={'small'}>
          <LogoutButton />
        </Box>
      </Header>
    )
  );
};

export default withAuth0(Navbar);
