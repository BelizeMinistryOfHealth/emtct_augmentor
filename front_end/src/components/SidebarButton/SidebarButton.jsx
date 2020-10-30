import React from 'react';
import { Button, Text, Box } from 'grommet';

const SidebarButton = ({ label, ...rest }) => {
  return (
    <Button plain {...rest}>
      {({ hover }) => (
        <Box
          background={hover ? 'accent-2' : undefined}
          pad={{ horizontal: 'medium', vertical: 'small' }}
        >
          <Text size='medium'>{label}</Text>
        </Box>
      )}
    </Button>
  );
};

export default SidebarButton;
