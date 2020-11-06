import React from 'react';
import { Box, Card, CardBody, Text } from 'grommet';
import { AddCircle, SubtractCircle, User } from 'grommet-icons';
import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';

const Identifier = ({ children, basicInfo, nextOfKin }) => {
  const {
    firstName,
    lastName,
    dob,
    patientId,
    ssn,
    countryOfBirth,
    address,
    education,
    ethnicity,
    hiv,
    hivDiagnosisDate,
  } = basicInfo;
  return (
    <Box gap='medium' align='center'>
      {children}
      <Box direction={'row'} gap={'medium'}>
        <Box>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            Name:
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            DOB:
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            Patient Id:
          </Text>
          <Text
            size={'small'}
            textAlign={'start'}
            gap={'medium'}
            weight={'bold'}
          >
            SSN:
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            Place of Birth:
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            Address:
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            Education:
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            Ethnicity:
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            HIV
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            HIV Diagnosis Date
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            Next of Kin:
          </Text>
          <Text size={'small'} textAlign={'start'} weight={'bold'}>
            Next of Kin Phone:
          </Text>
        </Box>

        <Box>
          <Text size={'small'} weight={'bold'} textAlign={'start'}>
            {firstName} {lastName}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {format(parseISO(dob), 'dd LLL yyyy')}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {patientId}
          </Text>

          <Text size={'small'} textAlign={'start'} gap={'medium'}>
            {ssn}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {countryOfBirth}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {address}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {education}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {ethnicity}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {hiv ? (
              <AddCircle color={'red'} />
            ) : (
              <SubtractCircle color={'green'} />
            )}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {format(parseISO(hivDiagnosisDate), 'dd LLL yyyy')}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {nextOfKin.name}
          </Text>
          <Text size={'small'} textAlign={'start'}>
            {nextOfKin.phoneNumber}
          </Text>
        </Box>
      </Box>
    </Box>
  );
};

const PatientBasicInfo = (props) => {
  const { basicInfo, nextOfKin } = props;

  return (
    <Card fill={'horizontal'}>
      <CardBody gap={'medium'} pad={'medium'}>
        <Identifier basicInfo={basicInfo} nextOfKin={nextOfKin}>
          <User size={'large'} />
        </Identifier>
      </CardBody>
    </Card>
  );
};

export default PatientBasicInfo;
