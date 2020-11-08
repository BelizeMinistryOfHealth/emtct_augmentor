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
import { Add, Close } from 'grommet-icons';
import React from 'react';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import { parseISO } from 'date-fns';
import format from 'date-fns/format';
import { Edit } from 'grommet-icons';
import EditForm from './HomeVisitEdit';
import { fetchHomeVisits } from '../../../api/patient';
import { useHttpApi } from '../../../providers/HttpProvider';
import { useHistory, useParams } from 'react-router-dom';
import AppCard from '../../AppCard/AppCard';

const homeVisitRow = (data, onClickEdit) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text size={'small'}>{data.reason}</Text>
      </TableCell>
      <TableCell align={'start'} alignContent={'start'}>
        <Text size={'small'} wordBreak={'break-word'} textAlign={'start'}>
          {data.comments}
        </Text>
      </TableCell>
      <TableCell align={'end'}>
        <Text size={'small'}>
          {format(parseISO(data.dateOfVisit), 'dd LLL yyyy')}
        </Text>
      </TableCell>

      <TableCell align={'end'}>
        <Text size={'small'}>
          {format(parseISO(data.createdAt), 'dd LLL yyyy')}
        </Text>
      </TableCell>
      <TableCell align={'end'}>
        <Text size={'small'}>{data.createdBy}</Text>
      </TableCell>
      <TableCell align={'center'} onClick={() => onClickEdit(data)}>
        <Edit size={'small'} />
      </TableCell>
    </TableRow>
  );
};

const HomeVisitsTable = ({ children, homeVisits, onClickEdit }) => {
  if (!homeVisits || homeVisits.length === 0) {
    return (
      <Box gap={'medium'} align={'center'}>
        <Text>Patient has not had any Home Visits.</Text>
      </Box>
    );
  }
  return (
    <Box gap={'medium'} align={'center'} width={'meidum'} fill={'horizontal'}>
      {children}
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'} colspan={2}>
              <Text size={'small'}>Reason</Text>
            </TableCell>
            <TableCell align={'start'} colspan={4}>
              <Text size={'small'}>Comments</Text>
            </TableCell>
            <TableCell align={'end'} colspan={1}>
              <Text size={'small'}>Date of Visit</Text>
            </TableCell>
            <TableCell colspan={1} align={'end'}>
              <Text size={'small'}>Date Created</Text>
            </TableCell>
            <TableCell colspan={2} align={'end'}>
              <Text size={'small'}>Created By</Text>
            </TableCell>
            <TableCell colspan={1} />
          </TableRow>
        </TableHeader>
        <TableBody>
          {homeVisits.map((h) => homeVisitRow(h, onClickEdit))}
        </TableBody>
      </Table>
    </Box>
  );
};

const HomeVisitList = (props) => {
  const { patientId } = useParams();
  const { httpInstance } = useHttpApi();
  const [editingHomeVisit, setEditingHomeVisit] = React.useState(undefined);
  const [data, setData] = React.useState({
    homeVisitsData: undefined,
    loading: false,
    error: undefined,
  });
  const [refreshHomeVisits, setRefreshHomeVisits] = React.useState(true);

  const history = useHistory();

  const onClickEdit = (homeVisit) => setEditingHomeVisit(homeVisit);

  React.useEffect(() => {
    const fetchVisits = async () => {
      try {
        setData({ homeVisitsData: undefined, loading: true, error: undefined });
        const homeVisits = await fetchHomeVisits(patientId, httpInstance);
        setData({
          homeVisitsData: homeVisits,
          loading: false,
          error: undefined,
        });
        setRefreshHomeVisits(false);
      } catch (e) {
        console.error(e);
        setData({
          homeVisitsData: undefined,
          loading: false,
          error: 'Failed to fetch home visits',
        });
      }
    };
    if (refreshHomeVisits) {
      fetchVisits();
    }
  }, [httpInstance, patientId, refreshHomeVisits]);

  const closeEditForm = () => {
    setEditingHomeVisit(undefined);
    setRefreshHomeVisits(true);
  };

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
                  Home Visits
                </Text>

                {data && data.homeVisitsData && (
                  <Text size={'large'} textAlign={'end'} weight={'normal'}>
                    {data.homeVisitsData.patient.firstName}{' '}
                    {data.homeVisitsData.patient.lastName}
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
                    history.push(`/patient/${patientId}/home_visits/new`)
                  }
                />
              </Box>
            </Box>
          </CardHeader>
          <CardBody gap={'medium'} pad={'medium'}>
            {editingHomeVisit && (
              <Layer
                position={'right'}
                full={'vertical'}
                onClickOutside={() => setEditingHomeVisit(undefined)}
                onEsc={() => setEditingHomeVisit(undefined)}
                modal
              >
                <Box
                  as={'form'}
                  fill={'vertical'}
                  overflow={'auto'}
                  width={'medium'}
                  pad={'medium'}
                >
                  <Box flex={false} direction={'row'} justify={'between'}>
                    <Heading level={2} margin={'none'}>
                      Edit
                    </Heading>
                    <Button
                      icon={<Close />}
                      onClick={() => setEditingHomeVisit(undefined)}
                    />
                  </Box>
                  <EditForm
                    visit={editingHomeVisit}
                    closeForm={closeEditForm}
                  />
                </Box>
              </Layer>
            )}

            {refreshHomeVisits ? (
              <Text>Loading....</Text>
            ) : (
              <HomeVisitsTable
                homeVisits={data.homeVisitsData.homeVisits}
                onClickEdit={onClickEdit}
              />
            )}
          </CardBody>
        </AppCard>
      </ErrorBoundary>
    </Layout>
  );
};

export default HomeVisitList;
