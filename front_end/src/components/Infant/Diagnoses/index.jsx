import { format, parseISO } from 'date-fns';
import {
  Box,
  CardBody,
  CardHeader,
  Heading,
  Table,
  TableCell,
  TableHeader,
  TableRow,
  TableBody,
  Text,
} from 'grommet';
import { InProgress } from 'grommet-icons';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import Layout from '../../Layout/Layout';
import InfantTabs from '../InfantTabs';

const diagnosisRow = (data) => {
  return (
    <TableRow key={data.diagnosisId}>
      <TableCell align={'start'}>
        <Text align={'start'} size={'small'}>
          {data.date ? format(parseISO(data.date), 'dd LLL yyyy') : 'N/A'}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'} align={'start'}>
          {data.name}
        </Text>
      </TableCell>
      <TableCell>
        <Text align={'start'} size={'small'}>
          {data.doctor}
        </Text>
      </TableCell>
      <TableCell>
        <Text align={'start'} size={'small'}>
          {data.comments}
        </Text>
      </TableCell>
    </TableRow>
  );
};

const DiagnosesTable = ({ children, data }) => {
  return (
    <Box gap={'medium'} pad={'medium'} align={'start'} fill>
      {children}
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text align={'start'}>Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Diagnosis</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Doctor</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Notes</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => diagnosisRow(d))}</TableBody>
      </Table>
    </Box>
  );
};

const InfantDiagnoses = () => {
  const { httpInstance } = useHttpApi();
  const { patientId, infantId, pregnancyId } = useParams();
  const [data, setData] = React.useState({
    result: undefined,
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const fetchData = async () => {
      try {
        console.log('fetching infant diagnoses');
        const result = await httpInstance.get(
          `/patients/${patientId}/infant/${infantId}/diagnoses`
        );
        setData({
          result: result.data,
          loading: false,
          error: undefined,
        });
      } catch (e) {
        console.error(e);
        setData({
          ...data,
          loading: false,
          error: e,
        });
      }
    };
    if (data.loading) {
      fetchData();
    }
  }, [patientId, httpInstance, data, infantId]);

  if (data.loading) {
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
          <Text>Loading... </Text>
          <InProgress size={'large'} />
        </Heading>
      </Box>
    );
  }

  if (data.error) {
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
        <InfantTabs data={data.result.infant} pregnancyId={pregnancyId}>
          <AppCard justify={'center'} gap={'medium'} fill={'horizontal'}>
            <CardHeader justify={'start'} pad={'medium'}>
              <Box>
                <span>
                  <Text size={'xxlarge'} weight={'bold'}>
                    Infant Diagnoses
                  </Text>
                  <span>
                    <Text size={'large'}>
                      {' '}
                      {data.result.infant.firstName}{' '}
                      {data.result.infant.lastName}
                    </Text>
                  </span>
                  <span>
                    <Text size={'medium'}>
                      {' '}
                      |{' '}
                      {format(parseISO(data.result.infant.dob), 'dd LLL yyyy')}
                    </Text>
                  </span>
                </span>

                <span>
                  <Text size={'medium'}>
                    <strong>Mother: </strong>
                  </Text>
                  <Text size={'medium'}>
                    {data.result.infant.mother.firstName}{' '}
                    {data.result.infant.mother.lastName}
                  </Text>
                </span>
              </Box>
            </CardHeader>
            <CardBody gap={'medium'} pad={'medium'}>
              <DiagnosesTable data={data.result.diagnoses} />
            </CardBody>
          </AppCard>
        </InfantTabs>
      </Box>
    </Layout>
  );
};

export default InfantDiagnoses;
