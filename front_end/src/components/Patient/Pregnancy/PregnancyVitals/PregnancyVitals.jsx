import { Box, CardBody, Text } from 'grommet';
import React from 'react';
import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';
import AppCard from '../../../AppCard/AppCard';

const VitalsComponent = ({ children, vitals, ...rest }) => {
  const { obstetricDetails } = vitals;
  return (
    <Box gap={'medium'} align={'center'} {...rest}>
      {children}
      <Box direction={'row'} gap={'medium'}>
        <Box>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            Para:
          </Text>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            C/S:
          </Text>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            Abortive Outcome:
          </Text>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            Date Pregnancy Diagnoses:
          </Text>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            Planned Pregnancy:
          </Text>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            Age at LMP:
          </Text>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            LMP:
          </Text>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            EDD:
          </Text>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            Interval b/w pregnancies:
          </Text>
        </Box>
        <Box>
          <Text size={'small'} textAlign={'start'}>
            {obstetricDetails.para}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {obstetricDetails.cs}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {obstetricDetails.pregnancyOutcome == ''
              ? 'N/A'
              : obstetricDetails.pregnancyOutcome}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {obstetricDetails.diagnosisDate
              ? format(parseISO(obstetricDetails.diagnosisDate), 'dd LLL yyyy')
              : 'N/A'}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {obstetricDetails.planned ? 'Yes' : 'No'}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {obstetricDetails.ageAtLmp}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {obstetricDetails.lmp
              ? format(parseISO(obstetricDetails.lmp), 'dd LLL yyyy')
              : 'N/A'}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {format(parseISO(obstetricDetails.edd), 'dd LLL yyyy')}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {vitals.interval}
          </Text>
        </Box>
      </Box>
    </Box>
  );
};

const PregnancyVitals = (props) => {
  return (
    <AppCard fill={'horizontal'}>
      <CardBody gap={'medium'} pad={'medium'}>
        <VitalsComponent vitals={props.vitals}>
          <Text size={'xlarge'} weight={'bold'}>
            Vitals
          </Text>
        </VitalsComponent>
      </CardBody>
    </AppCard>
  );
};

export default PregnancyVitals;
