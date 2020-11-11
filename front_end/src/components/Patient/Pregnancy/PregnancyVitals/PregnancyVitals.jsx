import { Box, CardBody, Text } from 'grommet';
import React from 'react';
import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';
import AppCard from '../../../AppCard/AppCard';

const VitalsComponent = ({ children, vitals, ...rest }) => {
  return (
    <Box gap={'medium'} align={'center'} {...rest}>
      {children}
      <Box direction={'row'} gap={'medium'}>
        <Box>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            Gestational Age:
          </Text>
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
            {vitals.gestationalAge}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {vitals.para}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {vitals.cs ? 'Yes' : 'No'}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {vitals.abortiveOutcome}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {vitals.diagnosisDate
              ? format(parseISO(vitals.diagnosisDate), 'dd LLL yyyy')
              : 'N/A'}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {vitals.planned ? 'Yes' : 'No'}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {vitals.ageAtLmp}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {vitals.lmp ? format(parseISO(vitals.lmp), 'dd LLL yyyy') : 'N/A'}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {format(parseISO(vitals.edd), 'dd LLL yyyy')}
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
  console.dir({ vitals: props.vitals });
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
