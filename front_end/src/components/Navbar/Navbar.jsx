import React from 'react';
import { Anchor, Box, Header } from 'grommet';
import Logout from '../Auth/Logout';

const Navbar = (props) => {
  const { permissions } = props;
  const isAdmin = permissions ? permissions.includes('admin:write') : false;
  return (
    <Header background={'dark-1'} pad={'small'}>
      <Box direction={'row'} align={'center'} gap={'small'}>
        <Anchor color={'white'} href={'/'}>
          EMTCT Project
        </Anchor>
        {isAdmin && (
          <>
            <Anchor color={'white'} href={'/admin/users'}>
              | Users
            </Anchor>
          </>
        )}
      </Box>
      <Box direction={'row'} align={'end'} gap={'small'}>
        <Logout />
      </Box>
    </Header>
  );
};

export default Navbar;
