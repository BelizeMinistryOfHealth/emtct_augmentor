import React from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import { Button } from 'grommet';

const LogoutButton = () => {
  const { isAuthenticated, logout } = useAuth0();

  return isAuthenticated && <Button label='Sign Out' onClick={logout} />;
};

export default LogoutButton;
