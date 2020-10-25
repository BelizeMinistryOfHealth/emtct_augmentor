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
import { Close } from 'grommet-icons';
import React from 'react';
import ErrorBoundary from '../../ErrorBoundary';
import Layout from '../../Layout/Layout';
import { parseISO } from 'date-fns';
import format from 'date-fns/format';
import { Edit } from 'grommet-icons';
import EditForm from './HomeVisitEdit';
import { fetchHomeVisits } from '../../../api/patient';
import { useHttpApi } from '../../../providers/HttpProvider';

const homeVisitRow = (data, onClickEdit) => {
  return (
    <TableRow key={data.id}>
      <TableCell align={'start'}>
        <Text>{data.reason}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.comments}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.dateOfVisit), 'dd LLL yyy')}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{data.createdBy}</Text>
      </TableCell>
      <TableCell align={'start'}>
        <Text>{format(parseISO(data.createdAt), 'dd LLL yyyy')}</Text>
      </TableCell>
      <TableCell align={'start'} onClick={() => onClickEdit(data)}>
        <Edit />
      </TableCell>
    </TableRow>
  );
};

const HomeVisitsTable = ({
  children,
  homeVisits,
  caption,
  onClickEdit,
  ...rest
}) => {
  if (homeVisits.length == 0) {
    return (
      <Box gap={'medium'} align={'center'}>
        <Text>Patient has not had any Home Visits.</Text>
      </Box>
    );
  }
  return (
    <Box gap={'medium'} align={'center'} width={'meidum'} fill={'horizontal'}>
      {children}
      <Table caption={caption}>
        <TableHeader>
          <TableRow>
            <TableCell align={'start'}>
              <Text>Reason</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Comments</Text>
            </TableCell>
            <TableCell align={'start'}>
              <Text>Date of Visit</Text>
            </TableCell>
            <TableCell>
              <Text>Created By</Text>
            </TableCell>
            <TableCell>
              <Text>Date Created</Text>
            </TableCell>
            <TableCell />
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
  const patientId = props.location.state.id;
  const { httpInstance } = useHttpApi();
  const [editingHomeVisit, setEditingHomeVisit] = React.useState(undefined);
  const onClickEdit = (homeVisit) => setEditingHomeVisit(homeVisit);
  const [data, setData] = React.useState({
    homeVisits: [],
    loading: false,
    error: undefined,
  });

  React.useEffect(() => {
    const fetchVisits = async () => {
      try {
        setData({ homeVisits: [], loading: true, error: undefined });
        const homeVisits = await fetchHomeVisits(patientId, httpInstance);
        setData({ homeVisits, loading: false, error: undefined });
      } catch (e) {
        setData({
          homeVisits: [],
          loading: false,
          error: 'Failed to fetch home visits',
        });
      }
    };
    fetchVisits();
  }, []);

  return (
    <Layout location={props.location} {...props}>
      <ErrorBoundary>
        <Card>
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
                  <EditForm visit={editingHomeVisit} />
                </Box>
              </Layer>
            )}
            <Box
              direction={'row-reverse'}
              align={'start'}
              pad={'medium'}
              gap={'medium'}
              background={'neutral-0'}
            >
              <Box align={'center'} pad={'medium'}>
                <Button label={'Create Home Visit'} />
              </Box>
            </Box>
            <HomeVisitsTable
              homeVisits={data.homeVisits}
              caption={'Home Visits'}
              onClickEdit={onClickEdit}
            />
          </CardBody>
        </Card>
      </ErrorBoundary>
    </Layout>
  );
};

export default HomeVisitList;
