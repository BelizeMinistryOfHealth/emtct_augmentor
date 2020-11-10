import { Box, CardBody, Text } from 'grommet';
import React from 'react';
import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';
import AppCard from '../../../AppCard/AppCard';

const CareInfo = ({ children, info, ...rest }) => {
  return (
    <Box gap={'medium'} align={'center'} {...rest}>
      {children}
      <Box direction={'row'} gap={'medium'}>
        <Box>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Date of Booking:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Gestation Age at Booking:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Prenatal Care Provider:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Total Prenatal Checks
          </Text>
        </Box>
        <Box>
          <Text size={'medium'} textAlign={'start'}>
            {info.dateOfBooking
              ? format(parseISO(info.dateOfBooking), 'dd LLL yyyy')
              : 'N/A'}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {info.gestationAge}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {info.prenatalCareProvider}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {info.totalChecks}
          </Text>
        </Box>
      </Box>
    </Box>
  );
};

const PreNatalCare = (props) => {
  return (
    <AppCard>
      <CardBody gap={'medium'} pad={'medium'}>
        <CareInfo info={props.info}>
          <Text size={'xlarge'} weight={'bold'}>
            Prenatal Care
          </Text>
        </CareInfo>
      </CardBody>
    </AppCard>
  );
};

export default PreNatalCare;
