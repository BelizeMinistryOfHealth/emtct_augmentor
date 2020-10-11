import React from 'react';
import { Box, Card, CardBody, Text } from 'grommet';
import { AddCircle, SubtractCircle, User } from 'grommet-icons';

const Identifier = ({ children, basicInfo, nextOfKin, ...rest }) => {
  const {
    firstName,
    lastName,
    dob,
    community,
    district,
    patientId,
    ssn,
    countryOfBirth,
    address,
    education,
    ethnicity,
    hiv,
  } = basicInfo;
  return (
    <Box gap='medium' align='center' {...rest}>
      {children}
      <Box direction={'row'} gap={'medium'}>
        <Box>
          <Text size={'medium'} weight={'bold'} textAlign={'start'}>
            Name:
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            DOB:
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            Patient Id:
          </Text>
          <Text
            size={'medium'}
            textAlign={'start'}
            gap={'medium'}
            weight={'bold'}
          >
            SSN:
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            Country of Birth:
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            Address:
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            Education:
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            Ethnicity:
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            HIV
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            Next of Kin:
          </Text>
          <Text size={'medium'} textAlign={'start'} weight={'bold'}>
            Next of Kin Phone:
          </Text>
        </Box>

        <Box>
          <Text size={'medium'} weight='bold' textAlign={'start'}>
            {firstName} {lastName}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {dob}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {patientId}
          </Text>

          <Text size={'medium'} textAlign={'start'} gap={'medium'}>
            {ssn}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {countryOfBirth}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {address}, {community}, {district}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {education}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {ethnicity}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {hiv ? (
              <AddCircle color={'red'} />
            ) : (
              <SubtractCircle color={'green'} />
            )}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
            {nextOfKin.name}
          </Text>
          <Text size={'medium'} textAlign={'start'}>
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
    <Card>
      <CardBody gap={'medium'} pad={'medium'}>
        <Identifier basicInfo={basicInfo} nextOfKin={nextOfKin}>
          <User size={'large'} />
        </Identifier>
      </CardBody>
    </Card>
  );
};

export default PatientBasicInfo;
