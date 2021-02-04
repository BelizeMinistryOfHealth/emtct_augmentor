import React from 'react';
import { Anchor, Box, Header } from 'grommet';
import Logout from '../Auth/Logout';
import { NavLink } from 'react-router-dom';
import './Navbar.css';

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
            <NavLink color={'white'} to={'/admin/users'}>
              | Users
            </NavLink>
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
