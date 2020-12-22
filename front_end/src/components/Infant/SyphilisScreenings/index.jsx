import { format, parseISO } from 'date-fns';
import {
  Box,
  CardBody,
  CardHeader,
  Heading,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import Layout from '../../Layout/Layout';
import Spinner from '../../Spinner';
import InfantTabs from '../InfantTabs';

const screeningRow = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell>
        <Text size={'small'} align={'start'}>
          {data.testName}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'} align={'start'}>
          {data.result}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text align={'start'} size={'small'}>
          {data.screeningDate
            ? format(parseISO(data.screeningDate), 'dd LLL yyyy')
            : 'N/A'}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text align={'start'} size={'small'}>
          {data.dateResultReceived
            ? format(parseISO(data.dateResultReceived), 'dd LLL yyyy')
            : 'N/A'}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text align={'start'} size={'small'}>
          {data.dateSampleTaken
            ? format(parseISO(data.dateSampleTaken), 'dd LLL yyyy')
            : 'N/A'}
        </Text>
      </TableCell>
    </TableRow>
  );
};

const ScreeningsTable = ({ children, data }) => {
  if (!data) {
    return (
      <Box alignContent={'center'}>
        <Text>No screenings found.</Text>
      </Box>
    );
  }
  return (
    <Box gap={'medium'} pad={'medium'} align={'start'} fill>
      {children}
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text align={'start'}>Test Name</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Result</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Screening Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Date Result Received</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Date Sample Taken</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => screeningRow(d))}</TableBody>
      </Table>
    </Box>
  );
};

const InfantSyphilisScreenings = () => {
  const { patientId, infantId, pregnancyId } = useParams();
  const { httpInstance } = useHttpApi();
  const [screeningsData, setScreeningsData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getScreenings = () => {
      httpInstance
        .get(`/patients/${patientId}/infant/${infantId}/syphilisScreenings`)
        .then((r) => {
          setScreeningsData({
            data: r.data,
            loading: false,
            error: undefined,
          });
        })
        .catch((e) => {
          console.error(e);
          setScreeningsData({
            data: undefined,
            loading: false,
            error: e.toJSON(),
          });
        });
    };
    if (screeningsData.loading) {
      getScreenings();
    }
  }, [screeningsData, infantId, patientId, httpInstance]);

  if (screeningsData.loading) {
    return (
      <Layout>
        <Box
          direction={'column'}
          gap={'large'}
          pad={'large'}
          justify={'center'}
          align={'center'}
          fill
        >
          <Heading>
            <Text>Loading.... </Text>
            <Spinner />
          </Heading>
        </Box>
      </Layout>
    );
  }

  if (screeningsData.error) {
    return (
      <Box
        direction={'colomn'}
        gap={'large'}
        pad={'large'}
        justify={'center'}
        align={'center'}
        fill
      >
        <Heading>
          <Text>Ooops. An error occurred while loading the data. </Text>
        </Heading>
      </Box>
    );
  }

  return (
    <Layout>
      <Box
        direction={'column'}
        gap={'xxlarge'}
        pad={{ left: 'small', bottom: 'xxsmall' }}
        justify={'evenly'}
        align={'center'}
        fill
      >
        <InfantTabs data={screeningsData.data.infant} pregnancyId={pregnancyId}>
          <AppCard justify={'center'} gap={'medium'} fill={'horizontal'}>
            <CardHeader justify={'start'} pad={'medium'}>
              <Box>
                <span>
                  <Text size={'xxlarge'} weight={'bold'}>
                    Infant Syphilis Screening
                  </Text>
                  <span>
                    <Text size={'large'}>
                      {' '}
                      {screeningsData.data.infant.firstName}{' '}
                      {screeningsData.data.infant.lastName}
                    </Text>
                  </span>
                  <span>
                    <Text size={'medium'}>
                      {' '}
                      |{' '}
                      {format(
                        parseISO(screeningsData.data.infant.dob),
                        'dd LLL yyyy'
                      )}
                    </Text>
                  </span>
                </span>
                <span>
                  <Text size={'medium'}>
                    <strong>Mother: </strong>
                  </Text>
                  <Text size={'medium'}>
                    {screeningsData.data.infant.mother.firstName}{' '}
                    {screeningsData.data.infant.mother.lastName}
                  </Text>
                </span>
              </Box>
            </CardHeader>
            <CardBody gap={'medium'} pad={'medium'}>
              <ScreeningsTable
                data={screeningsData.data.screenings}
              ></ScreeningsTable>
            </CardBody>
          </AppCard>
        </InfantTabs>
      </Box>
    </Layout>
  );
};

export default InfantSyphilisScreenings;
