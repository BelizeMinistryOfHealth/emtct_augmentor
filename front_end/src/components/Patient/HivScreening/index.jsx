import { format, parseISO } from 'date-fns';
import {
  Box,
  Button,
  CardBody,
  Table,
  TableCell,
  TableHeader,
  TableRow,
  TableBody,
  Text,
  Layer,
  Heading,
  CardHeader,
} from 'grommet';
import { Add, Checkmark, Close, Edit, Subtract } from 'grommet-icons';
import React from 'react';
import { useHistory, useParams } from 'react-router-dom';
import { useHttpApi } from '../../../providers/HttpProvider';
import AppCard from '../../AppCard/AppCard';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import EditForm from './HivScreeningEdit';

const screeningRow = (data, onClickEdit) => {
  return (
    <TableRow key={data.id} color={'red'}>
      <TableCell align={'end'}>
        {data.timely ? (
          <Checkmark size={'medium'} color={'blue'} />
        ) : (
          <Subtract size={'medium'} color={'red'} />
        )}
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.testName}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.result}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'} wordBreak={'break-word'}>
          {data.sampleCode}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.destination}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {format(parseISO(data.dueDate), 'dd LLL yyyy')}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {format(parseISO(data.screeningDate), 'dd LLL yyyy')}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {format(parseISO(data.dateSampleTaken), 'dd LLL yyyy')}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {data.dateSampleReceivedAtHq
            ? format(parseISO(data.dateSampleReceivedAtHq), 'dd LLL yyyy')
            : 'N/A'}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {data.dateResultReceived
            ? format(parseISO(data.dateResultReceived), 'dd LLL yyyy')
            : 'N/A'}
        </Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text size={'small'}>
          {data.dateResultShared
            ? format(parseISO(data.dateResultShared), 'dd LLL yyyy')
            : 'N/A'}
        </Text>
      </TableCell>
      <TableCell align={'start'} onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
    </TableRow>
  );
};

const HivScreeningTable = ({ children, caption, screenings, onClickEdit }) => {
  if (!screenings || screenings.length === 0) {
    return (
      <Box
        gap={'medium'}
        align={'center'}
        fill={'horizontal'}
        justify={'center'}
        width={'xlarge'}
      >
        <Text size={'xlarge'}>
          No Hiv Screenings available for this patient!
        </Text>
      </Box>
    );
  }

  screenings.sort((a, b) => a.dueDate > b.dueDate);

  return (
    <Box gap={'medium'} align={'center'} width={'medium'} fill={'horizontal'}>
      {children}
      <Table caption={caption}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'} />
            <TableCell align={'start'}>
              <Text size={'small'}>Test Name</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Test Result</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Sample Code</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Destination</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Due Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Screening Date</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Date Sample Taken</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Date Sample Received at HQ</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Date Result Received</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text size={'small'}>Date Result Shared</Text>{' '}
            </TableCell>
            <TableCell />
          </TableRow>
        </TableHeader>
        <TableBody>
          {screenings.map((i) => screeningRow(i, onClickEdit))}
        </TableBody>
      </Table>
    </Box>
  );
};

const HivScreening = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [data, setData] = React.useState({
    screenings: undefined,
    loading: true,
    error: undefined,
  });
  const [editingScreening, setEditingScreening] = React.useState(undefined);
  const history = useHistory();

  const onClickEdit = (screening) => setEditingScreening(screening);

  const closeEditScreen = () => {
    setEditingScreening(undefined);
    setData({ screenings: undefined, loading: true, error: undefined });
  };

  React.useEffect(() => {
    const fetchScreenings = async () => {
      try {
        // setData({ screenings: undefined, loading: true, error: undefined });
        const result = await httpInstance.get(
          `/patient/${patientId}/hivScreenings`
        );
        setData({ screenings: result.data, loading: false, error: undefined });
      } catch (e) {
        setData({
          screenings: undefined,
          loading: false,
          error: 'Failed to fetch hiv screenings',
        });
      }
    };
    if (data.loading) {
      fetchScreenings();
    }
  }, [httpInstance, patientId, data]);

  if (data.loading) {
    return <>Loading....</>;
  }

  if (data.error) {
    return <>Could not fetch patient HIV Screenings Data!</>;
  }

  return (
    <Layout location={props.location} {...props}>
      <ErrorBoundary>
        <AppCard fill={'horizontal'} pad={'small'} overflow={'scroll'}>
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
                  HIV Screenings
                </Text>
                {data && data.screenings.patient && (
                  <Text size={'large'} textAlign={'end'} weight={'normal'}>
                    {data.screenings.patient.firstName}{' '}
                    {data.screenings.patient.lastName}
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
                    history.push(`/patient/${patientId}/hiv_screenings/new`)
                  }
                />
              </Box>
            </Box>
          </CardHeader>
          <CardBody gap={'medium'} pad={'medium'}>
            {editingScreening && (
              <Layer
                position={'right'}
                full={'vertical'}
                onClickOutside={() => setEditingScreening(undefined)}
                onEsc={() => setEditingScreening(undefined)}
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
                      onClick={() => setEditingScreening(undefined)}
                    />
                  </Box>
                  <EditForm
                    screening={editingScreening}
                    closeEditScreen={closeEditScreen}
                  />
                </Box>
              </Layer>
            )}
            <HivScreeningTable
              screenings={data.screenings.hivScreenings}
              caption={'HIV Screenings'}
              onClickEdit={onClickEdit}
            />
          </CardBody>
        </AppCard>
      </ErrorBoundary>
    </Layout>
  );
};

export default HivScreening;
