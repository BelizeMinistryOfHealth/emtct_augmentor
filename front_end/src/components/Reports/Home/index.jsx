import React from 'react';
import { Box, Card, CardBody, CardFooter, Grid, Grommet, Text } from 'grommet';
import { Search, Tasks } from 'grommet-icons';
import { useHistory } from 'react-router-dom';

const data = [
  {
    color: 'blue',
    icon: <Search size={'large'} />,
    title: 'Search Patients',
    subTitle: 'Search pregnant women by BHIS patient id',
    message: 'Search for pregnant women in the BHIS',
    link: '/search',
  },
  {
    color: 'orange',
    icon: <Tasks size={'large'} />,
    title: 'Infant PCRs',
    subTitle: 'Infants with pending PCR Tests',
    message: 'PCRs',
    link: '/reports/pcrs',
  },
  // {
  //   color: 'green',
  //   icon: <System size={'large'} />,
  //   title: 'Summary PMTCT',
  //   subTitle: 'PMTCT Summary By Year',
  //   message: 'PMTCTs',
  //   link: '/',
  // },
];

const theme = {
  global: {
    font: {
      family: `-apple-system,
         BlinkMacSystemFont,
         "Segoe UI"`,
    },
    colors: {
      blue: '#00C8FF',
      green: '#17EBA0',
      teal: '#82FFF2',
      purple: '#F740FF',
      red: '#FC6161',
      orange: '#FFBC44',
      yellow: '#FFEB59',
    },
  },
  card: {
    footer: {
      pad: { horizontal: 'medium', vertical: 'small' },
      background: '#FFFFFF27',
    },
  },
};

const Identifier = ({ children, title, subTitle, size, ...rest }) => (
  <Box gap='small' align='center' {...rest}>
    {children}
    <Box>
      <Text size={size} weight='bold'>
        {title}
      </Text>
      <Text size={size}>{subTitle}</Text>
    </Box>
  </Box>
);

const ReportHome = () => {
  const history = useHistory();
  return (
    <Grommet theme={theme} full>
      <Box pad={'large'} margin={{ left: 'large', right: 'large' }}>
        <Grid
          gap={'medium'}
          rows={'small'}
          columns={{ count: 'fit', size: 'small' }}
          responisve={true}
        >
          {data.map((value) => (
            <Card
              background={value.color}
              key={value.message}
              flex={'shrink'}
              onClick={() => history.push(value.link)}
            >
              <CardBody pad={'small'}>
                <Identifier
                  pad={'small'}
                  title={value.title}
                  subTitle={value.subTitle}
                  size={'small'}
                  align={'start'}
                >
                  {value.icon}
                </Identifier>
              </CardBody>
              <CardFooter pad={{ horizontal: 'medium', vertical: 'small' }}>
                <Text size={'xsmall'}>{value.message}</Text>
              </CardFooter>
            </Card>
          ))}
        </Grid>
      </Box>
    </Grommet>
  );
};

export default ReportHome;
