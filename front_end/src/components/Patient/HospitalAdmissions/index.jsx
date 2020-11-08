import { format, parseISO } from 'date-fns';
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
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import EditForm from './HospitalAdmissionsEdit';
import AppCard from '../../AppCard/AppCard';

const admissionsRow = (data, onClickEdit) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {format(parseISO(data.dateAdmitted), 'dd LLL yyyy')}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.facility}</Text>
      </TableCell>
      <TableCell align={'center'} onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
    </TableRow>
  );
};

const AdmissionsTable = ({ children, caption, admissions, onClickEdit }) => {
  if (!admissions || admissions.length === 0) {
    return (
      <Box
        gap={'medium'}
        align={'center'}
        fill={'horizontal'}
        justify={'center'}
        width={'xlarge'}
      >
        <Text size={'xlarge'}>No Admissions were found for this patient!</Text>
      </Box>
    );
  }

  return (
    <Box gap={'medium'} align={'center'} width={'medium'} fill={'horizontal'}>
      {children}
      <Table caption={caption}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text size={'small'}>Date Admitted</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Facility</Text>
            </TableCell>
            <TableCell align={'start'} />
          </TableRow>
        </TableHeader>
        <TableBody>
          {admissions.map((i) => admissionsRow(i, onClickEdit))}
        </TableBody>
      </Table>
    </Box>
  );
};

const HospitalAdmissions = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [editingAdmission, setEditingAdmission] = React.useState(undefined);

  const [data, setData] = React.useState({
    admissionsData: undefined,
    loading: false,
    error: undefined,
  });
  const [refreshAdmissions, setRefreshAdmissions] = React.useState(true);
  const history = useHistory();

  const onClickEdit = (admission) => setEditingAdmission(admission);

  React.useEffect(() => {
    const fetchAdmissions = async () => {
      try {
        setData({ admissions: [], loading: true, error: undefined });
        const result = await httpInstance.get(
          `/patient/${patientId}/hospitalAdmissions`
        );
        setData({
          admissionsData: result.data,
          loading: false,
          error: undefined,
        });
        setRefreshAdmissions(false);
      } catch (e) {
        console.error(e);
        setData({
          admissions: undefined,
          loading: false,
          error: 'Failed to fetch hospital admissions!',
        });
      }
    };
    if (refreshAdmissions) {
      fetchAdmissions();
    }
  }, [httpInstance, patientId, refreshAdmissions]);

  const closeEditForm = () => {
    setEditingAdmission(undefined);
    setRefreshAdmissions(true);
  };

  if (refreshAdmissions) {
    return <>Loading...</>;
  }

  if (data.error) {
    return <>Could not fetch patient Hospital Admissions Data!</>;
  }

  return (
    <Layout location={props.location} {...props}>
      <ErrorBoundary>
        <AppCard fill={'horizontal'} pad={'small'}>
          <CardHeader>
            <Box direction={'row'} align={'start'} fill='horizontal'>
              <Box
                direction={'column'}
                align={'start'}
                fill={'horizontal'}
                justify={'between'}
                alignContent={'center'}
              >
                <Text size={'xxlarge'} weight={'bold'} textAlign={'start'}>
                  Hospital Admissions
                </Text>
                {data && data.admissionsData.patient && (
                  <Text size={'large'} textAlign={'end'} weight={'normal'}>
                    {data.admissionsData.patient.firstName}{' '}
                    {data.admissionsData.patient.lastName}
                  </Text>
                )}
              </Box>
              <Box
                align={'start'}
                fill={'horizontal'}
                direction={'row-reverse'}
              >
                <Button
                  icon={<Add />}
                  label={'Add'}
                  onClick={() =>
                    history.push(`/patient/${patientId}/admissions/new`)
                  }
                />
              </Box>
            </Box>
          </CardHeader>
          <CardBody gap={'medium'} pad={'medium'}>
            {editingAdmission && (
              <Layer
                position={'right'}
                full={'vertical'}
                onClickOutside={() => setEditingAdmission(undefined)}
                onEsc={() => setEditingAdmission(undefined)}
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
                    <Button
                      icon={<Close />}
                      onClick={() => setEditingAdmission(undefined)}
                    />
                  </Box>
                  <EditForm
                    admission={editingAdmission}
                    closeForm={closeEditForm}
                  />
                </Box>
              </Layer>
            )}
            <AdmissionsTable
              admissions={data.admissionsData.hospitalAdmissions}
              caption={'Hospital Admissions'}
              onClickEdit={onClickEdit}
            />
          </CardBody>
        </AppCard>
      </ErrorBoundary>
    </Layout>
  );
};

export default HospitalAdmissions;
