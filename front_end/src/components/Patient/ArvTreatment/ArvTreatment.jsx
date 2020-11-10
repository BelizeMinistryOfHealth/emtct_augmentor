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
import { InProgress } from 'grommet-icons';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import Layout from '../../Layout/Layout';

const row = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{data.pharmaceutical}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.strength}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.prescribedTime), 'dd LLL yyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.comments}</Text>
      </TableCell>
    </TableRow>
  );
};

const ArvsTable = ({ children, data, ...rest }) => {
  const arvs = data.arvs ?? [];

  return (
    <Box gap={'medium'} align={'center'} {...rest}>
      {children}
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text>ARV</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Dosage</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Comments</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{arvs.map((d) => row(d))}</TableBody>
      </Table>
    </Box>
  );
};

const BasicInfo = (props) => {
  const { basicInfo } = props;

  return (
    <Box direction={'row'} align={'start'} fill='horizontal'>
      <Box
        direction={'column'}
        align={'start'}
        fill={'horizontal'}
        justify={'between'}
        alignContent={'center'}
      >
        <Text size={'xxlarge'} weight={'bold'} textAlign={'start'}>
          ARVs
        </Text>
        <Text size={'large'} textAlign={'end'} weight={'normal'}>
          {basicInfo.firstName} {basicInfo.lastName}{' '}
        </Text>
      </Box>
    </Box>
  );
};

const ArvTreatment = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [arvData, setArvData] = React.useState({
    arvs: [],
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getArvs = async () => {
      try {
        const result = await httpInstance.get(`/patient/${patientId}/arvs`);
        setArvData({ arvs: result.data, loading: false, error: undefined });
      } catch (e) {
        console.error(e);
        setArvData({ arvs: [], loading: false, error: e });
      }
    };
    if (arvData.loading) {
      getArvs();
    }
  }, [httpInstance, arvData, patientId]);

  if (arvData.loading) {
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

  if (arvData.error) {
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
          <Text>Oooops. An error occurred while loading the data.</Text>
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
        <AppCard fill={'horizontal'}>
          <CardHeader gap={'medium'} pad={'medium'}>
            <BasicInfo basicInfo={arvData.arvs.patient} />
          </CardHeader>
          <CardBody gap={'medium'} pad={'medium'}>
            <ArvsTable data={arvData.arvs}></ArvsTable>
          </CardBody>
        </AppCard>
      </Box>
    </Layout>
  );
};

export default ArvTreatment;
