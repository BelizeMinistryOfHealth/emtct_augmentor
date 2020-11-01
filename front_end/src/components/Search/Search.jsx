import React from 'react';
import { withAuth0 } from '@auth0/auth0-react';
import { Box } from 'grommet';
import PatientIdSearch from './PatientIdSearch';
import { useHistory } from 'react-router-dom';

const SearchForm = () => {
  const history = useHistory();
  // Create handler for retrieving patient by patient id.
  const patientIdSearchHandler = (patientId) => {
    history.push(`/patient/${patientId}`, { id: patientId });
  };
  return (
    <Box
      fill
      align={'center'}
      justify={'center'}
      direction={'row-responsive'}
      pad={'large'}
      gap={'medium'}
    >
      <PatientIdSearch onSubmit={patientIdSearchHandler} />
    </Box>
  );
};

export default withAuth0(SearchForm);
