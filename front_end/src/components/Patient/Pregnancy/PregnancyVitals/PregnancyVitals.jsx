import { Box, Card, CardBody, Text } from 'grommet';
import { Checkmark, SubtractCircle } from 'grommet-icons';
import React from 'react';
import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';

const VitalsComponent = ({ children, vitals, ...rest }) => {
  return (
    <Box gap={'medium'} align={'center'} {...rest}>
      {children}
      <Box direction={'row'} gap={'medium'}>
        <Box>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Gestational Age:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Para:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            C/S:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Abortive Outcome:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Date Pregnancy Diagnoses:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Planned Pregnancy:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Age at LMP:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            LMP:
          </Text>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            EDD:
          </Text>
        </Box>
        <Box>
          <Text size={'medium'} textAlign={'start'}>
            {vitals.gestationalAge}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {vitals.para}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {vitals.cs ? <Checkmark /> : <SubtractCircle color={'green'} />}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {vitals.abortiveOutcome}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {format(parseISO(vitals.diagnosisDate), 'dd LLL yyyy')}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {vitals.planned ? (
              <Checkmark color={'green'} />
            ) : (
              <SubtractCircle color={'red'} />
            )}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {vitals.ageAtLmp}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {format(parseISO(vitals.LMP), 'dd LLL yyyy')}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {format(parseISO(vitals.EDD), 'dd LLL yyyy')}
          </Text>
        </Box>
      </Box>
    </Box>
  );
};

const PregnancyVitals = (props) => {
  return (
    <Card>
      <CardBody gap={'medium'} pad={'medium'}>
        <VitalsComponent vitals={props.vitals} />
      </CardBody>
    </Card>
  );
};

export default PregnancyVitals;
