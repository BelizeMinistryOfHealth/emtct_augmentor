import { format, parseISO } from 'date-fns';
import {
  Box,
  Button,
  Card,
  CardBody,
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
      <TableCell align={'start'} onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
    </TableRow>
  );
};

const AdmissionsTable = ({ children, caption, admissions, onClickEdit }) => {
  if (admissions.length === 0) {
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
    admissions: [],
    loading: false,
    error: undefined,
  });
  const history = useHistory();

  const onClickEdit = (admission) => setEditingAdmission(admission);

  React.useEffect(() => {
    const fetchAdmissions = async () => {
      try {
        setData({ admissions: [], loading: true, error: undefined });
        const result = await httpInstance.get(
          `/patient/${patientId}/hospitalAdmissions`
        );
        const admissions = result.data ?? [];
        setData({ admissions, loading: false, error: undefined });
      } catch (e) {
        setData({
          admissions: [],
          loading: false,
          error: 'Failed to fetch hospital admissions!',
        });
      }
    };
    fetchAdmissions();
  }, [httpInstance, patientId]);

  if (data.loading) {
    return <>Loading...</>;
  }

  if (data.error) {
    return <>Could not fetch patient Hospital Admissions Data!</>;
  }

  return (
    <Layout location={props.location} {...props}>
      <ErrorBoundary>
        <Card fill={'horizontal'}>
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
                  <EditForm admission={editingAdmission} />
                </Box>
              </Layer>
            )}
            <Box
              direction={'row-reverse'}
              align={'start'}
              pad={'medium'}
              gap={'medium'}
            >
              <Box align={'end'} pad={'medium'} fill={'horizontal'}>
                <Button
                  icon={<Add />}
                  label={'Add Hospital Admission'}
                  onClick={() =>
                    history.push(`/patient/${patientId}/admissions/new`)
                  }
                />
              </Box>
            </Box>
            <AdmissionsTable
              admissions={data.admissions}
              caption={'Hospital Admissions'}
              onClickEdit={onClickEdit}
            />
          </CardBody>
        </Card>
      </ErrorBoundary>
    </Layout>
  );
};

export default HospitalAdmissions;
