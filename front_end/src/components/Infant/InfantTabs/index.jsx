import { Box, Tab, Tabs, ThemeContext } from 'grommet';
import React from 'react';
import { useHistory } from 'react-router-dom';

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
        <Tab title={'Info'}>{content}</Tab>
        <Tab title={'HIV Screenings'}></Tab>
        <Tab title={'Diagnoses'}></Tab>
        <Tab title={'Syphilis Treatments'}></Tab>
        <Tab title={'Syphilis Screenings'}></Tab>
      </Tabs>
    </Box>
  );
};

const InfantTabs = ({ data, pregnancyId, children }) => {
  const { mother } = data;
  const history = useHistory();
  const onClickTab = (nextIndex) => {
    switch (nextIndex) {
      case 0:
        history.push(
          `/patient/${mother.patientId}/pregnancy/${pregnancyId}/infant/${data.id}`
        );
        break;
      case 1:
        history.push(
          `/patient/${mother.patientId}/infant/${data.id}/hivScreenings`
        );
        break;
      case 2:
        history.push(
          `/patient/${mother.patientId}/pregnancy/${pregnancyId}/infant/${data.id}/diagnoses`
        );
        break;
      case 3:
        history.push(
          `/patient/${mother.patientId}/pregnancy/${pregnancyId}/infant/${data.id}/syphilisTreatment`
        );
        break;
      case 4:
        history.push(
          `/patient/${mother.patientId}/pregnancy/${pregnancyId}/infant/${data.id}/syphilisScreenings`
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

export default InfantTabs;
