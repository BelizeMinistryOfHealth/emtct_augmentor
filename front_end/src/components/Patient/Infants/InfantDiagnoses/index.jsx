import { format, parseISO } from 'date-fns';
import {
  Box,
  Card,
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
import { InProgress } from 'grommet-icons';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useHttpApi } from '../../../../providers/HttpProvider';
import Layout from '../../../Layout/Layout';

const diagnosisRow = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text align={'start'}>
          {data.date ? format(parseISO(data.date), 'dd LLL yyyy') : 'N/A'}
        </Text>
      </TableCell>
      <TableCell>
        <Text align={'start'}>{data.diagnosis}</Text>
      </TableCell>
      <TableCell>
        <Text align={'start'}>{data.doctor}</Text>
      </TableCell>
      <TableCell>
        <Text align={'start'}>{data.comments}</Text>
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
  const { patientId } = useParams();
  const [data, setData] = React.useState({
    result: { diagnoses: [], patient: {} },
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const fetchData = async () => {
      try {
        const result = await httpInstance.get(
          `/patient/${patientId}/infants/diagnoses`
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
  }, [patientId, httpInstance, data]);

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
        pad={'medium'}
        justify={'evenly'}
        align={'center'}
        fill
      >
        <Card justify={'center'} gap={'medium'} fill={'horizontal'}>
          <CardHeader justify={'start'} pad={'medium'}>
            <Box>
              <Heading pad={'large'} gap={'medium'}>
                Infant Diagnoses
              </Heading>
              <span>
                <Text size={'large'}>
                  <strong>Mother: </strong>
                </Text>
                <Text size={'large'}>
                  {data.result.patient.firstName} {data.result.patient.lastName}
                </Text>
              </span>
            </Box>
          </CardHeader>
          <CardBody gap={'medium'} pad={'medium'}>
            <DiagnosesTable data={data.result.diagnoses} />
          </CardBody>
        </Card>
      </Box>
    </Layout>
  );
};

export default InfantDiagnoses;
