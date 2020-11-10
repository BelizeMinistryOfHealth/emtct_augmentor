import parseISO from 'date-fns/parseISO';
import format from 'date-fns/format';
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
import { fetchPregnancyLabResults } from '../../../api/patient';
import Layout from '../../Layout/Layout';
import { InProgress } from 'grommet-icons';
import AppCard from '../../AppCard/AppCard';

const labResultsRow = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text align={'start'}>
          {data.dateSampleTaken
            ? format(parseISO(data.dateSampleTaken), 'dd LLL yyyy')
            : 'N/A'}
        </Text>
      </TableCell>
      <TableCell>
        <Text align={'start'}>{data.testResult}</Text>
      </TableCell>
      <TableCell>
        <Text align={'start'} weight={'bold'}>
          {data.testName}
        </Text>
      </TableCell>
    </TableRow>
  );
};

const LabResultsTable = ({ children, data, ...rest }) => {
  const labTests = data ?? [];
  return (
    <Box
      gap={'medium'}
      pad={'medium'}
      align={'start'}
      fill={'horizontal'}
      {...rest}
    >
      {children}
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text align={'start'}> Date sample taken</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}> Test Result</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Test Name</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{labTests.map((d) => labResultsRow(d))}</TableBody>
      </Table>
    </Box>
  );
};

const LabResults = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [labData, setLabData] = React.useState({
    results: undefined,
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getLabResults = async () => {
      try {
        const results = await fetchPregnancyLabResults(patientId, httpInstance);
        setLabData({ results, loading: false, error: undefined });
      } catch (e) {
        console.error(e);
        setLabData({ results: undefined, loading: false, error: e });
      }
    };
    if (labData.loading) {
      getLabResults();
    }
  }, [labData, httpInstance, patientId]);

  if (labData.loading) {
    return (
      <Box
        direction={'column'}
        gap={'large'}
        pad={'large'}
        justify={'center'}
        align={'center'}
        fill
      >
        <Heading>
          <Text>Loading </Text>
          <InProgress />
        </Heading>
      </Box>
    );
  }

  if (labData.error) {
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
    <Layout props={props}>
      <Box
        direction={'column'}
        gap={'medium'}
        pad={'medium'}
        justify={'evenly'}
        align={'center'}
        fill
      >
        <AppCard fill={'horizontal'} pad={'small'}>
          <CardHeader justify={'evenly'}>
            <Box direction={'row'} align={'start'} fill='horizontal'>
              <Box
                direction={'column'}
                align={'start'}
                fill={'horizontal'}
                justify={'between'}
                alignContent={'center'}
              >
                <Text size={'xxlarge'} weight={'bold'} textAlign={'start'}>
                  Lab Tests
                </Text>
                <Text size={'large'} textAlign={'end'} weight={'normal'}>
                  {labData.results.patient.firstName}{' '}
                  {labData.results.patient.lastName}
                </Text>
              </Box>
            </Box>
          </CardHeader>
          <CardBody gap={'medium'} pad={'medium'} alignContent={'start'}>
            <LabResultsTable data={labData.results.labResults} />
          </CardBody>
        </AppCard>
      </Box>
    </Layout>
  );
};

export default LabResults;
