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

const AppTabs = ({ basicInfo, arvs, labResults }) => {
  return (
    <Tabs justify={'start'} flex>
      <Tab title={<RichTabTitle label={'Basic Info'} />}>{basicInfo}</Tab>
      <Tab title={<RichTabTitle label={'ARVs'} />}>{arvs}</Tab>
      <Tab title={<RichTabTitle label={'Lab Results'} />}> {labResults}</Tab>
    </Tabs>
  );
};

AppTabs.propTypes = {
  basicInfo: PropTypes.node.isRequired,
  arvs: PropTypes.node.isRequired,
  labResults: PropTypes.node.isRequired,
};

export default AppTabs;
