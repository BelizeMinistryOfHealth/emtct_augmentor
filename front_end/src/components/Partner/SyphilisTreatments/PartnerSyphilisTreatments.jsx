import { parseISO, format } from 'date-fns';
import {
  Box,
  Button,
  CardBody,
  CardHeader,
  Heading,
  Layer,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import { Add, Close, Edit } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import Layout from '../../Layout/Layout';
import Spinner from '../../Spinner';
import PartnerTabs from '../PartnerTabs/PartnerTabs';
import PartnerSyphilisTreatmentEdit from './PartnerSyphilisTreatmentEdit';

const treatmentRow = (data, onClickEdit) => {
  return (
    <TableRow key={data.id}>
      <TableCell>
        <Text size={'small'} align={'start'}>
          {data.medication}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'} align={'start'}>
          {data.dosage}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'} align={'start'}>
          {data.comments}
        </Text>
      </TableCell>
      <TableCell>
        <Text size={'small'} align={'start'}>
          {format(parseISO(data.date), 'dd LLL yyyy')}
        </Text>
      </TableCell>
      <TableCell onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
    </TableRow>
  );
};

const TreatmentsTable = ({ children, data, onClickEdit }) => {
  if (!data) {
    return (
      <Box alignContent={'center'}>
        <Text>No data found.</Text>
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
              <Text align={'start'}>Medication</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Dosage</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Comments</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text align={'start'}>Date</Text>
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>{data.map((d) => treatmentRow(d, onClickEdit))}</TableBody>
      </Table>
    </Box>
  );
};

const PartnerSyphilisTreatments = () => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [treatments, setTreatments] = React.useState({
    data: undefined,
    loading: true,
    error: undefined,
  });
  const [editingTreatment, setEditingTreatment] = React.useState(undefined);
  const onClickEdit = (treatment) => setEditingTreatment(treatment);
  const onCloseEditScreen = () => {
    setEditingTreatment(undefined);
    setTreatments({ data: undefined, loading: true, error: undefined });
  };
  const history = useHistory();

  React.useEffect(() => {
    const getTreatments = () => {
      httpInstance
        .get(`/patient/${patientId}/partners/syphilisTreatments`)
        .then((r) => {
          setTreatments({ data: r.data, loading: false, error: undefined });
        })
        .catch((e) => {
          console.error(e);
          setTreatments({ data: undefined, loading: false, error: e.toJSON() });
        });
    };
    if (treatments.loading) {
      getTreatments();
    }
  }, [patientId, treatments, httpInstance]);

  if (treatments.loading) {
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
            <Text>Loading....</Text>
            <Spinner />
          </Heading>
        </Box>
      </Layout>
    );
  }

  if (treatments.error) {
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
        <PartnerTabs data={treatments.data.patient}>
          <AppCard overflow={'scroll'} pad={'small'} fill={'horizontal'}>
            <CardHeader justify={'start'} pad={'medium'}>
              <Box direction={'row'} align={'start'} fill={'horizontal'}>
                <Box fill={'horizontal'}>
                  <span>
                    <Text size={'xxlarge'} weight={'bold'}>
                      Partner Syphilis Treatments
                    </Text>
                  </span>
                  <Text size={'large'}>
                    {' '}
                    {treatments.data.patient.firstName}{' '}
                    {treatments.data.patient.lastName}
                  </Text>
                </Box>
                <Box
                  align={'start'}
                  fill={'horizontal'}
                  direction={'row-reverse'}
                >
                  <Button
                    icon={<Add />}
                    label={'add'}
                    onClick={() =>
                      history.push(
                        `/patient/${patientId}/partners/syphilisTreatments/new`
                      )
                    }
                  />
                </Box>
              </Box>
            </CardHeader>
            <CardBody gap={'medium'} pad={'medium'}>
              {editingTreatment && (
                <Layer
                  position={'right'}
                  full={'vertical'}
                  onClickOutside={onCloseEditScreen}
                  modal
                >
                  <Box
                    fill={'vertical'}
                    overflow={'auto'}
                    width={'medium'}
                    pad={'medium'}
                  >
                    <Box
                      flex={false}
                      direction={'row-responsive'}
                      justify={'between'}
                    >
                      <Heading level={2} margin={'none'}>
                        Edit
                      </Heading>
                      <Button icon={<Close />} onClick={onCloseEditScreen} />
                    </Box>
                    <PartnerSyphilisTreatmentEdit
                      treatment={editingTreatment}
                      closeEditScreen={onCloseEditScreen}
                    />
                  </Box>
                </Layer>
              )}
              <TreatmentsTable
                data={treatments.data.treatments}
                onClickEdit={onClickEdit}
              />
            </CardBody>
          </AppCard>
        </PartnerTabs>
      </Box>
    </Layout>
  );
};

export default PartnerSyphilisTreatments;
