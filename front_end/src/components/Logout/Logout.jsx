import React from 'react';
import { Button } from 'grommet';

const LogoutButton = () => {
  return <Button label='Sign Out' onClick={() => console.log('logout')} />;
};

export default LogoutButton;
