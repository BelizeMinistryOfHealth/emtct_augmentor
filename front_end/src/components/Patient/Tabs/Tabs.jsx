import React from 'react';
import PropTypes from 'prop-types';
import { Tab, Tabs, Box, Text } from 'grommet';

const RichTabTitle = ({ icon, label }) => (
  <Box direction='row' align='center' gap='xsmall' margin='xsmall'>
    {icon}
    <Text size='small'>
      <strong>{label}</strong>
    </Text>
  </Box>
);

RichTabTitle.propTypes = {
  label: PropTypes.string.isRequired,
};

const AppTabs = ({ basicInfo }) => {
  return (
    <Tabs justify={'start'} flex={'grow'}>
      <Tab title={<RichTabTitle label={'Basic Info'} />}>{basicInfo}</Tab>
    </Tabs>
  );
};

AppTabs.propTypes = {
  basicInfo: PropTypes.node.isRequired,
};

export default AppTabs;
