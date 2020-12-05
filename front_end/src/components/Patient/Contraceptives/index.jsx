import { format, parseISO } from 'date-fns';
import {
  Box,
  Button,
  CardBody,
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
  Text,
  Layer,
  Heading,
  CardHeader,
} from 'grommet';
import { Add, Close, Edit } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import EditForm from './ContraceptivesEdit';

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
  if (!contraceptives || contraceptives.length === 0) {
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
    result: undefined,
    loading: true,
    error: undefined,
  });
  const [editingContraceptive, setEditingContraceptive] = React.useState(
    undefined
  );
  const history = useHistory();

  const onClickEdit = (contraceptive) => setEditingContraceptive(contraceptive);

  const onCloseForm = () => {
    setData({ result: undefined, loading: true, error: undefined });
    setEditingContraceptive(undefined);
  };

  React.useEffect(() => {
    const fetchContraceptives = async () => {
      try {
        const result = await httpInstance.get(
          `/patients/${patientId}/contraceptivesUsed`
        );
        setData({ result: result.data, loading: false, error: undefined });
      } catch (e) {
        setData({
          result: undefined,
          loading: false,
          error: 'Error occurred while fetching contraceptive information!',
        });
      }
    };
    if (data.loading) {
      fetchContraceptives();
    }
  }, [httpInstance, patientId, data]);
  if (data.loading) {
    return <>Loading....</>;
  }

  return (
    <Layout {...props}>
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
                  Contraceptives Used
                </Text>
                {data && data.result && data.result.patient && (
                  <Text size={'large'} textAlign={'end'} weight={'normal'}>
                    {data.result.patient.firstName}{' '}
                    {data.result.patient.lastName}
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
                    history.push(`/patient/${patientId}/contraceptives/new`)
                  }
                />
              </Box>
            </Box>
          </CardHeader>
          <CardBody gap={'medium'} pad={'medium'}>
            {editingContraceptive && (
              <Layer
                position={'right'}
                full={'vertical'}
                onClickOutside={() => setEditingContraceptive(undefined)}
                onEsc={() => setEditingContraceptive(undefined)}
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
                      onClick={() => setEditingContraceptive(undefined)}
                    />
                  </Box>
                  <EditForm
                    contraceptive={editingContraceptive}
                    onCloseForm={onCloseForm}
                  />
                </Box>
              </Layer>
            )}
            <ContraceptivesTable
              contraceptives={
                data && data.result ? data.result.contraceptives : []
              }
              caption={'Contraceptives Used'}
              onClickEdit={onClickEdit}
            />
          </CardBody>
        </AppCard>
      </ErrorBoundary>
    </Layout>
  );
};

export default ContraceptivesUsed;
