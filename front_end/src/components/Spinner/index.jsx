import React from 'react';
import { Box } from 'grommet';
import { ReactComponent as SpinnerSvg } from './spinner.svg';

const Spinner = (props) => {
  let size = props.size;
  if (!size) {
    size = 228;
  }
  return (
    <Box align='center' justify='center'>
      <SpinnerSvg height={size} width={size} />
    </Box>
  );
};

export default Spinner;
