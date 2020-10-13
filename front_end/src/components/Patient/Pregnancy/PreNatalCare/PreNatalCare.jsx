import { Box, Card, CardBody, Text } from 'grommet';
import React from 'react';
import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';

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
        </Box>
        <Box>
          <Text size={'medium'} textAlign={'start'}>
            {format(parseISO(info.dateOfBooking), 'dd LLL yyyy')}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {info.gestationAge}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {info.prenatalCareProvider}
          </Text>
        </Box>
      </Box>
    </Box>
  );
};

const PreNatalCare = (props) => {
  return (
    <Card>
      <CardBody gap={'medium'} pad={'medium'}>
        <CareInfo info={props.info}>
          <Text size={'xlarge'} weight={'bold'}>
            Prenatal Care
          </Text>
        </CareInfo>
      </CardBody>
    </Card>
  );
};

export default PreNatalCare;
