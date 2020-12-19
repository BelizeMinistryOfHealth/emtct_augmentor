import React from 'react';
import 'firebase/auth';
import firebase from 'firebase';
import { Button } from 'grommet';

const Logout = () => {
  // Log out function
  const handleClick = () => {
    firebase.auth().signOut();
  };

  return <Button label='Sign Out' onClick={handleClick} />;
};
export default Logout;
