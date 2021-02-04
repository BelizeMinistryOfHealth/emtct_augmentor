import React from 'react';
import { Box, Text } from 'grommet';
import SearchField from './SearchField';

const Search = ({ onSubmit, label, errMessage }) => {
  return (
    <Box
      fill
      align={'center'}
      justify={'center'}
      direction={'column'}
      pad={'small'}
      gap={'small'}
    >
      {errMessage.length > 0 && (
        <Box>
          <Text color={'accent-4'}>{errMessage}</Text>
        </Box>
      )}
      <SearchField onSubmit={onSubmit} label={label} />
    </Box>
  );
};

export default Search;
