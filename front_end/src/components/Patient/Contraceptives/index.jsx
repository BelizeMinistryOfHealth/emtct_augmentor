import { format, parseISO } from 'date-fns';
import {
  Box,
  Button,
  CardBody,
  Card,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
} from 'grommet';
import { Add, Edit } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';

const contraceptiveRow = (data, onClickEdit) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{data.contraceptive}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.dateUsed), 'dd LLL yyyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.comments}</Text>
      </TableCell>
      <TableCell align={'start'} onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
    </TableRow>
  );
};

const ContraceptivesTable = ({
  children,
  contraceptives,
  caption,
  onClickEdit,
}) => {
  if (contraceptives.length === 0) {
    return (
      <Box gap={'medium'} align={'center'}>
        <Text>No Contraceptives Information Exists for this Patient!</Text>
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
              <Text>Contraceptive</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Date Used</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Comments</Text>
            </TableCell>
            <TableCell align={'start'}></TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>
          {contraceptives.map((i) => contraceptiveRow(i, onClickEdit))}
        </TableBody>
      </Table>
    </Box>
  );
};

const ContraceptivesUsed = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [data, setData] = React.useState({
    contraceptives: [],
    loading: false,
    error: undefined,
  });
  const [editingContraceptive, setEditingContraceptive] = React.useState(
    undefined
  );
  const history = useHistory();

  const onClickEdit = (contraceptive) => setEditingContraceptive(contraceptive);

  React.useEffect(() => {
    const fetchContraceptives = async () => {
      try {
        setData({ contraceptives: [], loading: true, error: undefined });
        const result = await httpInstance.get(
          `/patient/${patientId}/contraceptivesUsed`
        );
        const contraceptives = result.data ?? [];
        setData({ contraceptives, loading: false, error: undefined });
      } catch (e) {
        setData({
          contraceptives: [],
          loading: false,
          error: 'Error occurred while fetching contraceptive information!',
        });
      }
    };
    fetchContraceptives();
  }, [httpInstance, patientId]);
  if (data.loading) {
    return <>Loading....</>;
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
              <Box align={'end'} pad={'medium'} fill={'horizontal'}>
                <Button
                  icon={<Add />}
                  label={'Add Contraceptive'}
                  onClick={() =>
                    history.push(`/patient/${patientId}/contraceptives/new`)
                  }
                />
              </Box>
            </Box>
            <ContraceptivesTable
              contraceptives={data.contraceptives}
              caption={'Contraceptives Used'}
              onClickEdit={onClickEdit}
            />
          </CardBody>
        </Card>
      </ErrorBoundary>
    </Layout>
  );
};

export default ContraceptivesUsed;
