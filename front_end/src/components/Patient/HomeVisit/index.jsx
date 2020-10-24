import {
  Box,
  Button,
  Card,
  CardBody,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import { useRecoilValueLoadable } from 'recoil';
import React from 'react';
import { homeVisitsSelector } from '../../../state';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import { parseISO } from 'date-fns';
import format from 'date-fns/format';

const homeVisitRow = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{data.testName}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.testResult}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.dateSampleTaken), 'dd LLL yyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.resultDate), 'dd LLL yyyy')}</Text>
      </TableCell>
    </TableRow>
  );
};

const HomeVisitsTable = ({ children, homeVisits, caption, ...rest }) => {
  if (homeVisits.length == 0) {
    return (
      <Box gap={'medium'} align={'center'}>
        <Text>Patient has not had any Home Visits.</Text>
      </Box>
    );
  }
  return (
    <Box gap={'medium'} align={'center'}>
      {children}
      <Table caption={caption}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text>Test Name</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Test Result</Text>
            </TableCell>
            <TableCell>
              <Text>Date Sample Taken</Text>
            </TableCell>
            <TableCell>
              <Text>Result Date</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{homeVisits.map((h) => homeVisitRow(h))}</TableBody>
      </Table>
    </Box>
  );
};

const HomeVisitList = (props) => {
  const patientId = props.location.state.id;
  const { state, contents } = useRecoilValueLoadable(
    homeVisitsSelector(patientId)
  );
  let homeVisits = [];
  switch (state) {
    case 'hasValue':
      homeVisits = contents;
      break;
    case 'hasError':
      return contents.message;
    case 'loading':
      return 'Loading....';
    default:
      return 'An unexpected error happened.';
  }

  return (
    <Layout location={props.location} {...props}>
      <ErrorBoundary>
        <Card>
          <CardBody gap={'medium'} pad={'medium'}>
            <Box
              direction={'row-reverse'}
              align={'start'}
              pad={'medium'}
              gap={'medium'}
              background={'neutral-0'}
            >
              <Box align={'center'} pad={'medium'}>
                <Button label={'Create Home Visit'} />
              </Box>
            </Box>
            <HomeVisitsTable homeVisits={homeVisits} />
          </CardBody>
        </Card>
      </ErrorBoundary>
    </Layout>
  );
};

export default HomeVisitList;
