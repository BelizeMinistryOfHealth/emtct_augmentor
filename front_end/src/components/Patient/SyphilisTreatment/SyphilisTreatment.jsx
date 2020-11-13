import {
  Box,
  CardBody,
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
import AppCardHeader from '../../AppCard/AppCardHeader';
import Layout from '../../Layout/Layout';
import { format, parseISO } from 'date-fns';

const row = (data) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.pharmaceutical}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.strength}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {format(parseISO(data.prescribedTime), 'dd LLL yyy')}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.comments}</Text>
      </TableCell>
    </TableRow>
  );
};

const PrescriptionsTable = ({ children, data, ...rest }) => {
  const prescriptions = data.prescriptions ?? [];

  return (
    <Box gap={'medium'} align={'center'} {...rest} fill={'horizontal'}>
      {children}
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text>Treatment</Text>
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
        <TableBody>{prescriptions.map((d) => row(d))}</TableBody>
      </Table>
    </Box>
  );
};

const SyphilisTreatment = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [treatmentData, setTreatmentData] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });

  React.useEffect(() => {
    const getTreatment = async () => {
      try {
        const result = await httpInstance.get(
          `/patient/${patientId}/syphilisTreatments`
        );
        setTreatmentData({
          data: result.data,
          loading: false,
          error: undefined,
        });
      } catch (e) {
        console.error(e);
        setTreatmentData({ data: undefined, loading: false, error: e });
      }
    };
    if (treatmentData.loading) {
      getTreatment();
    }
  }, [treatmentData, httpInstance, patientId]);

  if (treatmentData.loading) {
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
          <Text>Loading...</Text>
        </Box>
      </Layout>
    );
  }
  return (
    <Layout>
      <Box
        direction={'column'}
        gap={'medium'}
        pad={'medium'}
        justify={'evenly'}
        align={'center'}
        fill
      >
        <AppCard fill={'horizontal'}>
          <AppCardHeader
            gap={'medium'}
            pad={'medium'}
            title={'Syphilis Treatments'}
            patient={treatmentData.data.patient}
          />
          <CardBody gap={'medium'} pad={'medium'}>
            <PrescriptionsTable data={treatmentData.data}></PrescriptionsTable>
          </CardBody>
        </AppCard>
      </Box>
    </Layout>
  );
};

export default SyphilisTreatment;
