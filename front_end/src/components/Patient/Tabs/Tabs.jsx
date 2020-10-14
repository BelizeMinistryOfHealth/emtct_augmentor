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
  icon: PropTypes.node.isRequired,
  label: PropTypes.string.isRequired,
};

const AppTabs = ({ basicInfo, arvs }) => {
  return (
    <Tabs>
      <Tab title={<RichTabTitle label={'Basic Info'} />}>{basicInfo}</Tab>
      <Tab title={<RichTabTitle label={'ARVs'} />}>{arvs}</Tab>
    </Tabs>
  );
};

AppTabs.propTypes = {
  basicInfo: PropTypes.node.isRequired,
  arvs: PropTypes.node.isRequired,
};

export default AppTabs;
