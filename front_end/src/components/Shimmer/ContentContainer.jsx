import React from 'react';
import PropTypes from 'prop-types';
import { Box, Heading, Text } from 'grommet';

export const ContentContainer = ({ children, title, content, ...rest }) => (
  <Box
    round
    border={{ color: 'grey' }}
    pad='medium'
    gap='small'
    direction='column'
    background='white'
    {...rest}
  >
    <Heading level={2} margin='none' size='small'>
      {title}
    </Heading>
    <Text color='gray' size='small'>
      {content}
    </Text>
    {children}
  </Box>
);

ContentContainer.propTypes = {
  title: PropTypes.string,
  content: PropTypes.string,
  children: PropTypes.node.isRequired,
};

export default ContentContainer;
