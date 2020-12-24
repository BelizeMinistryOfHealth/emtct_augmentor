import { Box, Tab, Tabs } from 'grommet';
import React from 'react';
import { useHistory } from 'react-router-dom';
import { ThemeContext } from 'styled-components';

const colors = {
  border: '#999999',
  'border-strong': '#666666',
  'border-weak': '#BBBBBB',
  'active-background': 'background-contrast',
  'active-text': 'text',
};

const customTheme = {
  global: {
    colors,
  },
  tab: {
    border: {
      disabled: {
        color: 'border-weak',
      },
    },
    disabled: {
      color: 'text-weak',
    },
  },
};

const TabsItem = ({ content, onClickTab }) => {
  const [index, setIndex] = React.useState(0);
  const onActive = (nextIndex) => {
    setIndex(nextIndex);
    onClickTab(nextIndex);
  };
  return (
    <Box gap='medium' pad='medium' fill={'horizontal'} alignSelf={'start'}>
      <Tabs
        activeIndex={index}
        onActive={onActive}
        alignControls={'start'}
        alignContent={'start'}
        fill={'horizontal'}
      >
        <Tab title={'Syphilis Treatment'}>{content}</Tab>
        <Tab title={'Contact Tracing'}>{content}</Tab>
      </Tabs>
    </Box>
  );
};

const PartnerTabs = ({ data, pregnancyId, children }) => {
  const { patient } = data;
  const history = useHistory();
  const onClickTab = (nextIndex) => {
    switch (nextIndex) {
      case 0:
        history.push(
          `/patient/${patient.patientId}/pregnancy/${pregnancyId}/partners/syphilisTreatments`
        );
        break;
      case 1:
        history.push(
          `/patient/${patient.patientId}/pregnancy/${pregnancyId}/partners/contactTracing`
        );
        break;
    }
  };

  return (
    <Box
      fill={'horizontal'}
      pad={{ left: 'xxsmall', top: 'xxsmall' }}
      alignSelf={'start'}
      alignContent={'start'}
    >
      <ThemeContext.Extend value={customTheme}>
        <TabsItem onClickTab={onClickTab} content={children} />
      </ThemeContext.Extend>
    </Box>
  );
};

export default PartnerTabs;
