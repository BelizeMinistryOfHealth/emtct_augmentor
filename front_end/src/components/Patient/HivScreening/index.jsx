import { format, parseISO } from 'date-fns';
import {
  Box,
  Button,
  Card,
  CardBody,
  Table,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import { Add } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';

const screeningRow = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{data.testName}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.result}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.sampleCode}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.destination}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.screeningDate), 'dd LLL yyyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>
          {format(parseISO(data.dateSampleReceivedAtHq), 'dd LLL yyyy')}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.dateResultReceived), 'dd LLL yyyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.dateResultShared), 'dd LLL yyyy')}</Text>
      </TableCell>
    </TableRow>
  );
};

const HivScreeningTable = ({ children, caption, screenings }) => {
  if (screenings.length === 0) {
    return (
      <Box
        gap={'medium'}
        align={'center'}
        fill={'horizontal'}
        justify={'center'}
        width={'xlarge'}
      >
        <Text size={'xlarge'}>
          No Hiv Screenings available for this patient!
        </Text>
      </Box>
    );
  }

  return (
    <Box gap={'medium'} align={'center'} width={'medium'} fill={'horizontal'}>
      {children}
      <Table caption={caption}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>Test Name</TableCell>
            <TableCell align={'start'}>Test Result</TableCell>
            <TableCell align={'start'}>Sample Code</TableCell>
            <TableCell align={'start'}>Destination</TableCell>
            <TableCell align={'start'}>Screening Date</TableCell>
            <TableCell align={'start'}>Date Sample Received at HQ</TableCell>
            <TableCell align={'start'}>Date Result Received</TableCell>
            <TableCell align={'start'}>Date Result Shared</TableCell>
          </TableRow>
        </TableHeader>
      </Table>
    </Box>
  );
};

const HivScreening = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [data, setData] = React.useState({
    screenings: [],
    loading: false,
    error: undefined,
  });
  const history = useHistory();

  React.useEffect(() => {
    const fetchScreenings = async () => {
      try {
        setData({ screenings: [], loading: true, error: undefined });
        const result = await httpInstance.get(
          `/patient/${patientId}/hivScreenings`
        );
        const screenings = result.data ?? [];
        setData({ screenings, loading: false, error: undefined });
      } catch (e) {
        setData({
          screenings: [],
          loading: false,
          error: 'Failed to fetch hiv screenings',
        });
      }
    };
    fetchScreenings();
  }, [httpInstance, patientId]);

  if (data.loading) {
    return <>Loading....</>;
  }

  if (data.error) {
    return <>Could not fetch patient's HIV Screenings Data!</>;
  }

  return (
    <Layout location={props.location} {...props}>
      <ErrorBoundary>
        <Card fill={'horizontal'}>
          <CardBody gap={'medium'} pad={'medium'}>
            <Box
              direction={'row-reverse'}
              align={'start'}
              pad={'medium'}
              gap={'medium'}
            >
              <Box align={'center'} pad={'medium'} fill={'horizontal'}>
                <Button
                  icon={<Add />}
                  label={'Create Hiv Screening'}
                  onClick={() =>
                    history.push(`/patient/${patientId}/hiv_screenings/new`)
                  }
                />
              </Box>
            </Box>
            <HivScreeningTable
              screenings={data.screenings}
              caption={'HIV Screenings'}
            />
          </CardBody>
        </Card>
      </ErrorBoundary>
    </Layout>
  );
};

export default HivScreening;
